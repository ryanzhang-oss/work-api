package controllers

import (
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/dynamic/fake"
	workv1alpha1 "sigs.k8s.io/work-api/pkg/apis/v1alpha1"
	"testing"
)

var (
	testGvr = schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "Deployment",
	}
	testDeployment = v12.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
	}
	testPod               = v1.Pod{}
	testInvalidYaml       = []byte(getRandomString())
	rawTestDeployment, _  = json.Marshal(testDeployment)
	rawInvalidResource, _ = json.Marshal(testInvalidYaml)
	rawMissingResource, _ = json.Marshal(testPod)
	testManifest          = workv1alpha1.Manifest{RawExtension: runtime.RawExtension{
		Raw: rawTestDeployment,
	}}

	InvalidManifest = workv1alpha1.Manifest{RawExtension: runtime.RawExtension{
		Raw: rawInvalidResource,
	}}

	MissingManifest = workv1alpha1.Manifest{RawExtension: runtime.RawExtension{
		Raw: rawMissingResource,
	}}
)

// This interface is needed for testMapper abstract class.
type testMappingInterface interface {
	meta.RESTMapper
	RESTMapping(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error)
}

type testMapper struct {
	meta.RESTMapper
}

func (m testMapper) RESTMapping(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error) {
	if gk.Kind == "Deployment" {
		return &meta.RESTMapping{
			Resource:         testGvr,
			GroupVersionKind: testDeployment.GroupVersionKind(),
			Scope:            nil,
		}, nil
	} else {
		return nil, nil
	}

}

func TestDecodeUnstructured(t *testing.T) {

	expectedGvr := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "Deployment",
	}

	emptyGvr := schema.GroupVersionResource{}
	testCases := map[string]struct {
		reconciler ApplyWorkReconciler
		manifest   workv1alpha1.Manifest
		wantGvr    schema.GroupVersionResource
		wantErr    bool
	}{
		"manifest is in proper format/ happy path": {
			reconciler: ApplyWorkReconciler{
				client:             &test.MockClient{},
				spokeDynamicClient: &fake.FakeDynamicClient{},
				spokeClient:        &test.MockClient{},
				log:                logr.Logger{},
				restMapper:         testMapper{},
			},
			manifest: testManifest,
			wantGvr:  expectedGvr,
			wantErr:  false,
		},
		"manifest has incorrect syntax/ decode fail": {
			reconciler: ApplyWorkReconciler{
				client:             &test.MockClient{},
				spokeDynamicClient: &fake.FakeDynamicClient{},
				spokeClient:        &test.MockClient{},
				log:                logr.Logger{},
				restMapper:         testMapper{},
			},
			manifest: InvalidManifest,
			wantGvr:  emptyGvr,
			wantErr:  true,
		},
		"manifest is correct / object not mapped in restmapper / decode fail": {
			reconciler: ApplyWorkReconciler{
				client:             &test.MockClient{},
				spokeDynamicClient: &fake.FakeDynamicClient{},
				spokeClient:        &test.MockClient{},
				log:                logr.Logger{},
				restMapper:         testMapper{},
			},
			manifest: MissingManifest,
			wantGvr:  emptyGvr,
			wantErr:  true,
		},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			gvr, obj, err := testCase.reconciler.decodeUnstructured(testCase.manifest)
			assert.Equalf(t, testCase.wantErr, err != nil, "Testcase %s", testName)
			if obj != nil {
				assert.Equalf(t, testGvr.Group, obj.GroupVersionKind().Group, "Testcase %s", testName)
				assert.Equalf(t, testGvr.Version, obj.GroupVersionKind().Version, "Testcase %s", testName)
				assert.Equalf(t, testDeployment.Kind, obj.GroupVersionKind().Kind, "Testcase %s", testName)
			}
			assert.Equalf(t, testCase.wantGvr, gvr, "Testcase %s", testName)
		})
	}
}

func TestBuildResourceIdentifier(t *testing.T) {
	testIndex := utilrand.Int()
	testGroup := getRandomString()
	testVersion := getRandomString()
	testKind := getRandomString()
	testNamespace := getRandomString()
	testName := getRandomString()
	objMap := make(map[string]interface{})
	objMap["apiVersion"] = testGroup + "/" + testVersion
	objMap["group"] = testGroup
	objMap["kind"] = testKind
	metadataMap := make(map[string]interface{})
	metadataMap["name"] = testName
	metadataMap["namespace"] = testNamespace
	objMap["metadata"] = map[string]string{}
	objMap["metadata"] = metadataMap
	unstructuredObj := unstructured.Unstructured{
		Object: objMap,
	}
	gvr := schema.GroupVersionResource{
		Group:    getRandomString(),
		Version:  getRandomString(),
		Resource: getRandomString(),
	}

	t.Run("resource identifier created", func(t *testing.T) {
		resourceIdentifier := buildResourceIdentifier(testIndex, &unstructuredObj, gvr)
		assert.Equalf(t, testIndex, resourceIdentifier.Ordinal, "Testcase Index Incorrect")
		assert.Equalf(t, unstructuredObj.Object["group"], resourceIdentifier.Group, "Testcase Group Incorrect")
		assert.Equalf(t, testVersion, resourceIdentifier.Version, "Testcase Version Incorrect")
		assert.Equalf(t, testKind, resourceIdentifier.Kind, "Testcase Kind Incorrect")
		assert.Equalf(t, testNamespace, resourceIdentifier.Namespace, "Testcase Namespace Incorrect")
		assert.Equalf(t, testName, resourceIdentifier.Name, "Testcase Name Incorrect")
		assert.Equalf(t, gvr.Resource, resourceIdentifier.Resource, "Testcase Resource Incorrect")
	})
}

func getRandomString() string {
	return utilrand.String(10)
}
