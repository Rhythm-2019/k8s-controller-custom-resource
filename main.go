/*
Copyright The Kubernetes Authors.

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
    "flag"
    clientset "github.com/Rhythm-2019/k8s-controller-custom-resource/pkg/client/clientset/versioned"
    informers "github.com/Rhythm-2019/k8s-controller-custom-resource/pkg/client/informers/externalversions"
    "github.com/Rhythm-2019/k8s-controller-custom-resource/pkg/signals"
    "github.com/golang/glog"
    utilruntime "k8s.io/apimachinery/pkg/util/runtime"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "time"
)

var (
    MasterUrl string
    Kubeconfig string
)
func init() {
    flag.StringVar(&Kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
    flag.StringVar(&MasterUrl, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func main() {
    flag.Parse()

    stopCh := signals.SetupSignalHandler()

    cfg, err := clientcmd.BuildConfigFromFlags(MasterUrl, Kubeconfig)
    if err != nil {
        glog.Fatalf("Error to build config from flags, detail is %v", err)
    }

    kubeClientSet, err := kubernetes.NewForConfig(cfg)
    if err != nil {
        glog.Fatalf("Error to create kube client, detail is %v", err)
    }
    networkClientSet, err := clientset.NewForConfig(cfg)
    if err != nil {
        glog.Fatalf("Error to create network client, detail is %v", err)
    }


    informerFactory := informers.NewSharedInformerFactory(networkClientSet, 10*time.Second)

    controller := NewController(kubeClientSet, networkClientSet, informerFactory.Simplecrd().V1().Networks())

    go informerFactory.Start(stopCh)
    if err := controller.Run(2, stopCh); err != nil {
        utilruntime.HandleError(err)
    }
}


