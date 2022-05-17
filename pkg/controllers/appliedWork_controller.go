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
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	workv1alpha1 "sigs.k8s.io/work-api/pkg/apis/v1alpha1"
)

// AppliedWorkReconciler reconciles an AppliedWork object
type AppliedWorkReconciler struct {
	appliedResourceTracker
	clusterNameSpace string
}

func newAppliedWorkReconciler(clusterNameSpace string, hubClient client.Client, spokeClient client.Client,
	spokeDynamicClient dynamic.Interface, restMapper meta.RESTMapper) *AppliedWorkReconciler {
	return &AppliedWorkReconciler{
		appliedResourceTracker: appliedResourceTracker{
			hubClient:          hubClient,
			spokeClient:        spokeClient,
			spokeDynamicClient: spokeDynamicClient,
			restMapper:         restMapper,
		},
		clusterNameSpace: clusterNameSpace,
	}
}

// Reconcile implement the control loop logic for AppliedWork object.
func (r *AppliedWorkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.InfoS("applied work reconcile loop triggered", "item", req.NamespacedName)
	appliedWork := &workv1alpha1.AppliedWork{}
	appliedWorkDeleted := false
	err := r.spokeClient.Get(ctx, req.NamespacedName, appliedWork)
	switch {
	case errors.IsNotFound(err):
		klog.InfoS("appliedWork does not exist", "item", req.NamespacedName)
		appliedWork = nil
		appliedWorkDeleted = true
	case err != nil:
		klog.ErrorS(err, "failed to get appliedWork", "item", req.NamespacedName)
		return ctrl.Result{}, err
	default:
		klog.InfoS("get appliedWork in the member cluster", "item", req.NamespacedName)
	}
	nsWorkName := req.NamespacedName
	nsWorkName.Namespace = r.clusterNameSpace
	if _, err := r.reconcile(ctx, nil, appliedWork, nsWorkName); err != nil {
		return ctrl.Result{}, err
	}
	// stop the periodic check if it's gone
	if appliedWorkDeleted {
		return ctrl.Result{}, nil
	}
	// we want to periodically check if what we've applied matches what is recorded
	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

// SetupWithManager wires up the controller.
func (r *AppliedWorkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&workv1alpha1.AppliedWork{}).Complete(r)
}
