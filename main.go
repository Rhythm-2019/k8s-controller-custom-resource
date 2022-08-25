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
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"time"
)

var (
	MasterUrl  string
	Kubeconfig string
)

func init() {
	flag.StringVar(&Kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&MasterUrl, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func main() {
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	klog.Infof("controller boot...")

	// 从 InClusterConfig 中获取配置
	cfg, err := clientcmd.BuildConfigFromFlags(MasterUrl, Kubeconfig)
	if err != nil {
		klog.Fatalf("Error to build config from flags, detail is %v", err)
	}

	// 可用于创建获取 kubenetes 中其他 API 资源的 informer，比如监听 Pod
	kubeClientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error to create kube client, detail is %v", err)
	}
	// 创建 CRD 后可以通过 RestAPI 获取资源信息，networkClient 封装了这些操作
	networkClientSet, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error to create network client, detail is %v", err)
	}

	// informer 工厂，用于创建 informer，informer 是带有缓存、用于监听资源变化的客户端
	informerFactory := informers.NewSharedInformerFactoryWithOptions(networkClientSet, 10*time.Second, informers.WithNamespace("default"))
	// Conrtoller 通过控制循环获取资源变更信息，并将实际状态转换到期望状态
	controller := NewController(kubeClientSet, networkClientSet, informerFactory.Samplecrd().V1().Networks())

	go informerFactory.Start(stopCh)
	if err := controller.Run(2, stopCh); err != nil {
		utilruntime.HandleError(err)
	}
}
