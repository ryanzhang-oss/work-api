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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	workv1alpha1 "sigs.k8s.io/work-api/pkg/apis/v1alpha1"
	"sigs.k8s.io/work-api/pkg/client/clientset/versioned"
)

// FinalizeWorkReconciler reconciles a Work object for finalization
type FinalizeWorkReconciler struct {
	client      client.Client
	spokeClient *versioned.Clientset
	restMapper  meta.RESTMapper
	log         logr.Logger
}

// Reconcile implement the control loop logic for finalizing Work object.
func (r *FinalizeWorkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	work := &workv1alpha1.Work{}
	err := r.client.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, work)
	switch {
	case errors.IsNotFound(err):
		return ctrl.Result{}, nil
	case err != nil:
		return ctrl.Result{}, err
	}

	klog.InfoS("Finalize work reconcile loop triggered", "item", req.NamespacedName)

	// cleanup finalizer and resources
	if !work.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(work, workFinalizer) {
			deletePolicy := metav1.DeletePropagationForeground
			err := r.spokeClient.MulticlusterV1alpha1().AppliedWorks().Delete(ctx, req.Name,
				metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
			if err != nil {
				klog.ErrorS(err, "failed to delete the applied Work", req.NamespacedName.String())
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(work, workFinalizer)
		}
		return ctrl.Result{}, r.client.Update(ctx, work, &client.UpdateOptions{})
	}

	// don't add finalizer to instances that already have it
	if controllerutil.ContainsFinalizer(work, workFinalizer) {
		return ctrl.Result{}, nil
	}

	klog.InfoS("appliedWork does not exist yet, we will create it", "item", req.NamespacedName)
	appliedWork := &workv1alpha1.AppliedWork{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Spec: workv1alpha1.AppliedWorkSpec{
			ManifestWorkName: req.Name,
		},
	}
	appliedWork, err = r.spokeClient.MulticlusterV1alpha1().AppliedWorks().Create(ctx, appliedWork, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		// if this conflicts, we'll simply try again later
		klog.ErrorS(err, "failed to create the appliedWork", "name", req.Name)
		return ctrl.Result{}, err
	}

	work.Finalizers = append(work.Finalizers, workFinalizer)
	return ctrl.Result{}, r.client.Update(ctx, work, &client.UpdateOptions{})
}

// SetupWithManager wires up the controller.
func (r *FinalizeWorkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&workv1alpha1.Work{},
		builder.WithPredicates(predicate.GenerationChangedPredicate{})).Complete(r)
}
