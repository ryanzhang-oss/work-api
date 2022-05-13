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
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	workv1alpha1 "sigs.k8s.io/work-api/pkg/apis/v1alpha1"
)

// AppliedWorkReconciler reconciles an AppliedWork object
type AppliedWorkReconciler struct {
	hubClient   client.Client
	spokeClient client.Client
	restMapper  meta.RESTMapper
}

// Reconcile implement the control loop logic for AppliedWork object.
func (r *AppliedWorkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	work := &workv1alpha1.Work{}
	err := r.spokeClient.Get(ctx, req.NamespacedName, work)
	switch {
	case errors.IsNotFound(err):
		klog.InfoS("appliedWork does not exist", "item", req.NamespacedName)
		return ctrl.Result{}, nil
	case err != nil:
		klog.ErrorS(err, "failed to get appliedWork", "item", req.NamespacedName)
		return ctrl.Result{}, err
	}

	klog.InfoS("appliedWork does not exist", "item", req.NamespacedName)

	// do nothing if the finalizer is not present
	// it ensures all maintained resources will be cleaned once work is deleted
	if !controllerutil.ContainsFinalizer(work, workFinalizer) {
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager wires up the controller.
func (r *AppliedWorkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&workv1alpha1.AppliedWork{}).Complete(r)
}
