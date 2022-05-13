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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"sigs.k8s.io/work-api/pkg/apis/v1alpha1"
	"sigs.k8s.io/work-api/pkg/controllers"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var hubkubeconfig string
	var workNamespace string
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&hubkubeconfig, "hub-kubeconfig", "", "Paths to a kubeconfig connect to hub.")
	flag.StringVar(&workNamespace, "work-namespace", "", "Namespace to watch for work.")
	flag.Parse()
	opts := ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		Port:               9443,
		Namespace:          workNamespace,
	}
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	hubConfig, err := getKubeConfig(hubkubeconfig)
	if err != nil {
		setupLog.Error(err, "error reading kubeconfig to connect to hub")
		os.Exit(1)
	}

	if err := controllers.Start(ctrl.SetupSignalHandler(), hubConfig, ctrl.GetConfigOrDie(), setupLog, opts); err != nil {
		setupLog.Error(err, "problem running controllers")
		os.Exit(1)
	}
}

func getKubeConfig(hubkubeconfig string) (*restclient.Config, error) {
	spokeClientSet, err := kubernetes.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		return nil, errors.Wrap(err, "cannot create the spoke client")
	}

	secret, err := spokeClientSet.CoreV1().Secrets("work").Get(context.Background(), hubkubeconfig, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "cannot find kubeconfig secrete")
	}

	kubeConfigData, ok := secret.Data["kubeconfig"]
	if !ok || len(kubeConfigData) == 0 {
		return nil, fmt.Errorf("wrong formatted kube config")
	}

	kubeConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigData)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create the rest client")
	}

	return kubeConfig, nil
}
