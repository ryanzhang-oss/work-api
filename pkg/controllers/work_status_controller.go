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
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	workv1alpha1 "sigs.k8s.io/work-api/pkg/apis/v1alpha1"
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

	work := &workv1alpha1.Work{}
	// fetch work CR from the member cluster
	err := r.hubClient.Get(ctx, req.NamespacedName, work)
	switch {
	case errors.IsNotFound(err):
		klog.InfoS("work does not exist", "item", req.NamespacedName)
		work = nil
	case err != nil:
		klog.ErrorS(err, "failed to get work", "item", req.NamespacedName)
		return ctrl.Result{}, err
	default:
		klog.InfoS("work exists in the hub cluster", "item", req.NamespacedName)
	}

	return r.reconcile(ctx, work, nil, req.NamespacedName)
}

// SetupWithManager wires up the controller.
func (r *WorkStatusReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&workv1alpha1.Work{},
		builder.WithPredicates(UpdateDeleteOnlyPredicate{})).Complete(r)
}

// We don't need to process t
type UpdateDeleteOnlyPredicate struct {
	predicate.Funcs
}

func (UpdateDeleteOnlyPredicate) Create(e event.CreateEvent) bool {
	return false
}

func (UpdateDeleteOnlyPredicate) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil {
		klog.Error("Update event has no old object to update", "event", e)
		return false
	}
	if e.ObjectNew == nil {
		klog.Error("Update event has no new object to update", "event", e)
		return false
	}
	return e.ObjectNew.GetResourceVersion() != e.ObjectOld.GetResourceVersion()
}
