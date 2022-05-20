/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	workapi "sigs.k8s.io/work-api/pkg/apis/v1alpha1"
)

// WorkStatusReconciler reconciles a Work object when its status changes
type WorkStatusReconciler struct {
	appliedResourceTracker
}

func newWorkStatusReconciler(hubClient client.Client, spokeClient client.Client, spokeDynamicClient dynamic.Interface, restMapper meta.RESTMapper) *WorkStatusReconciler {
	return &WorkStatusReconciler{
		appliedResourceTracker{
			hubClient:          hubClient,
			spokeClient:        spokeClient,
			spokeDynamicClient: spokeDynamicClient,
			restMapper:         restMapper,
		},
	}
}

// Reconcile implement the control loop logic for Work Status.
func (r *WorkStatusReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.InfoS("work status reconcile loop triggered", "item", req.NamespacedName)
	work, appliedWork, err := r.fetchWorks(ctx, req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	// work has been garbage collected
	if work == nil {
		return ctrl.Result{}, nil
	}

	// from now on both work objects should exist
	newRes, staleRes := r.calculateNewAppliedWork(work, appliedWork)
	if err = r.deleteStaleWork(ctx, staleRes); err != nil {
		klog.ErrorS(err, "failed to delete all the stale work", "work", req.NamespacedName)
		// we can't proceed to update the applied
		return ctrl.Result{}, err
	}

	// update the appliedWork with the new work
	appliedWork.Status.AppliedResources = newRes
	if err = r.spokeClient.Status().Update(ctx, appliedWork, &client.UpdateOptions{}); err != nil {
		klog.ErrorS(err, "update appliedWork status failed", "appliedWork", appliedWork.GetName())
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// calculateNewAppliedWork check the difference between what is supposed to be applied  (tracked by the work CR status)
// and what was applied in the member cluster (tracked by the appliedWork CR).
// What is in the `appliedWork` but not in the `work` should be deleted from the member cluster
// What is in the `work` but not in the `appliedWork` should be added to the appliedWork status
func (r *WorkStatusReconciler) calculateNewAppliedWork(work *workapi.Work, appliedWork *workapi.AppliedWork) ([]workapi.AppliedManifestResourceMeta, []workapi.AppliedManifestResourceMeta) {
	var staleRes, newRes []workapi.AppliedManifestResourceMeta

	for _, resourceMeta := range appliedWork.Status.AppliedResources {
		resStillExist := false
		for _, manifestCond := range work.Status.ManifestConditions {
			if isSameResource(resourceMeta, manifestCond.Identifier) {
				resStillExist = true
				break
			}
		}
		if !resStillExist {
			klog.V(3).InfoS("find an orphaned resource", "parent work", work.GetObjectKind().GroupVersionKind(),
				"name", work.GetName(), "resource", resourceMeta)
			staleRes = append(staleRes, resourceMeta)
		}
	}

	for _, manifestCond := range work.Status.ManifestConditions {
		ac := meta.FindStatusCondition(manifestCond.Conditions, ConditionTypeApplied)
		if ac == nil {
			klog.Errorf("find one work %+v that has no applied condition", manifestCond.Identifier)
			continue
		}
		// we only add the applied one to the appliedWork status
		if ac.Status == metav1.ConditionTrue {
			resRecorded := false
			// we keep the existing resourceMeta since it has the UID
			for _, resourceMeta := range appliedWork.Status.AppliedResources {
				if isSameResource(resourceMeta, manifestCond.Identifier) {
					resRecorded = true
					newRes = append(newRes, resourceMeta)
					break
				}
			}
			if !resRecorded {
				klog.V(5).InfoS("find a new resource", "parent work", work.GetObjectKind().GroupVersionKind(),
					"name", work.GetName(), "resource", manifestCond.Identifier)
				newRes = append(newRes, workapi.AppliedManifestResourceMeta{
					ResourceIdentifier: manifestCond.Identifier,
				})
			}
		}
	}

	return newRes, staleRes
}

func (r *WorkStatusReconciler) deleteStaleWork(ctx context.Context, staleWorks []workapi.AppliedManifestResourceMeta) error {
	var errs []error

	for _, staleWork := range staleWorks {
		gvr := schema.GroupVersionResource{
			Group:    staleWork.Group,
			Version:  staleWork.Version,
			Resource: staleWork.Resource,
		}
		err := r.spokeDynamicClient.Resource(gvr).Namespace(staleWork.Namespace).
			Delete(ctx, staleWork.Name, metav1.DeleteOptions{})
		if err != nil && !errors.IsGone(err) {
			klog.ErrorS(err, "failed to delete a stale work", "work", staleWork)
			errs = append(errs, err)
		}

	}
	return utilerrors.NewAggregate(errs)
}

// isSameResource checks if an appliedMeta is referring to the same resource that a resourceId is pointing to
func isSameResource(appliedMeta workapi.AppliedManifestResourceMeta, resourceId workapi.ResourceIdentifier) bool {
	return appliedMeta.Resource == resourceId.Resource && appliedMeta.Version == resourceId.Version &&
		appliedMeta.Group == resourceId.Group && appliedMeta.Namespace == resourceId.Namespace &&
		appliedMeta.Name == resourceId.Name
}

// SetupWithManager wires up the controller.
func (r *WorkStatusReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&workapi.Work{},
		builder.WithPredicates(UpdateOnlyPredicate{}, predicate.ResourceVersionChangedPredicate{})).Complete(r)
}

// We only need to process the update event
type UpdateOnlyPredicate struct {
	predicate.Funcs
}

func (UpdateOnlyPredicate) Create(event.CreateEvent) bool {
	return false
}

func (UpdateOnlyPredicate) Delete(event.DeleteEvent) bool {
	return false
}
