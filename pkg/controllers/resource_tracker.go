package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	workv1alpha1 "sigs.k8s.io/work-api/pkg/apis/v1alpha1"
)

type appliedResourceTracker struct {
	hubClient          client.Client
	spokeClient        client.Client
	spokeDynamicClient dynamic.Interface
	restMapper         meta.RESTMapper
}

func (r *appliedResourceTracker) reconcile(ctx context.Context,
	work *workv1alpha1.Work, appliedWork *workv1alpha1.AppliedWork, nsWorkName types.NamespacedName) (ctrl.Result, error) {
	if work == nil {
		work = &workv1alpha1.Work{}
		// fetch work CR from the member cluster
		err := r.hubClient.Get(ctx, nsWorkName, work)
		switch {
		case errors.IsNotFound(err):
			klog.InfoS("work does not exist", "item", nsWorkName)
			work = nil
		case err != nil:
			klog.ErrorS(err, "failed to get work", "item", nsWorkName)
			return ctrl.Result{}, err
		default:
			klog.InfoS("work exists in the hub cluster", "item", nsWorkName)
		}
	}

	if appliedWork == nil {
		appliedWork = &workv1alpha1.AppliedWork{}
		// fetch appliedWork CR from the member cluster
		err := r.spokeClient.Get(ctx, nsWorkName, appliedWork)
		switch {
		case errors.IsNotFound(err):
			klog.InfoS("appliedWork does not exist", "item", nsWorkName)
			appliedWork = nil
		case err != nil:
			klog.ErrorS(err, "failed to get appliedWork", "item", nsWorkName)
			return ctrl.Result{}, err
		default:
			klog.InfoS("appliedWork exists in the member cluster", "item", nsWorkName)
		}
	}

	if err := checkConsistentExist(work, appliedWork, nsWorkName); err != nil {
		klog.ErrorS(err, "applied/work object existence not consistent", "item", nsWorkName)
		return ctrl.Result{}, err
	}

	if err := r.removeDeletedAppliedWork(ctx, work, appliedWork); err != nil {
		klog.ErrorS(err, "failed to calculate the difference between the work and what we have applied", nsWorkName)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// removeDeletedAppliedWork check the difference between what is supposed to be applied  (tracked by the work CR status)
// and what was applied in the member cluster (tracked by the appliedWork CR) and remove those are applied but no longer exist in the work
func (r *appliedResourceTracker) removeDeletedAppliedWork(ctx context.Context, work *workv1alpha1.Work, appliedWork *workv1alpha1.AppliedWork) error {
	if work == nil && appliedWork == nil {
		klog.InfoS("both applied and work are garbage collected")
		return nil
	}
	return nil
}

func checkConsistentExist(work *workv1alpha1.Work, appliedWork *workv1alpha1.AppliedWork, workName types.NamespacedName) error {
	// work already deleted
	if work == nil && appliedWork != nil {
		return fmt.Errorf("work finalizer didn't delete the appliedWork %s", workName)
	}
	// we are triggered by appliedWork change or work update so the appliedWork should already be here
	if work != nil && appliedWork == nil {
		return fmt.Errorf("work controller didn't create the appliedWork %s", workName)
	}
	return nil
}

func (r *appliedResourceTracker) decodeUnstructured(manifest workv1alpha1.Manifest) (schema.GroupVersionResource, *unstructured.Unstructured, error) {
	unstructuredObj := &unstructured.Unstructured{}
	err := unstructuredObj.UnmarshalJSON(manifest.Raw)
	if err != nil {
		return schema.GroupVersionResource{}, nil, fmt.Errorf("Failed to decode object: %w", err)
	}
	mapping, err := r.restMapper.RESTMapping(unstructuredObj.GroupVersionKind().GroupKind(), unstructuredObj.GroupVersionKind().Version)
	if err != nil {
		return schema.GroupVersionResource{}, nil, fmt.Errorf("Failed to find gvr from restmapping: %w", err)
	}

	return mapping.Resource, unstructuredObj, nil
}
