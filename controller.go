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
    "fmt"
    simplecrdv1 "github.com/Rhythm-2019/k8s-controller-custom-resource/pkg/apis/simplecrd/v1"
    clientset "github.com/Rhythm-2019/k8s-controller-custom-resource/pkg/client/clientset/versioned"
    networkscheme "github.com/Rhythm-2019/k8s-controller-custom-resource/pkg/client/clientset/versioned/scheme"
    v1 "github.com/Rhythm-2019/k8s-controller-custom-resource/pkg/client/informers/externalversions/simplecrd/v1"
    listers "github.com/Rhythm-2019/k8s-controller-custom-resource/pkg/client/listers/simplecrd/v1"
    "github.com/golang/glog"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    utilruntime "k8s.io/apimachinery/pkg/util/runtime"
    "k8s.io/apimachinery/pkg/util/wait"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/kubernetes/scheme"
    typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/tools/record"
    "k8s.io/client-go/util/workqueue"
    "time"
)

const (
    controllerComponentName = "network-controller"
    SuccessSynced = "Synced"
    MessageResourceSynced = "Network synced successfully"
)

type Controller struct {
    kubeClientSet   kubernetes.Interface
    networkClientSet  clientset.Interface
    networkLister    listers.NetworkLister
    netwrokSynced   cache.InformerSynced
    recorder    record.EventRecorder
    workQueue   workqueue.RateLimitingInterface

}
func NewController(kubeClientSet *kubernetes.Clientset, networkClientSet *clientset.Clientset, networkInformer v1.NetworkInformer) *Controller {

    utilruntime.Must(networkscheme.AddToScheme(scheme.Scheme))
    glog.V(4).Infof("creating event broadcaster")
    broadcaster := record.NewBroadcaster()
    broadcaster.StartLogging(glog.Infof)
    broadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientSet.CoreV1().Events("")})
    recorder := broadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerComponentName})

    controller := &Controller{
        kubeClientSet:    kubeClientSet,
        networkClientSet: networkClientSet,
        networkLister:    networkInformer.Lister(),
        netwrokSynced:    networkInformer.Informer().HasSynced,
        recorder:         recorder,
        workQueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "network"),
    }

    glog.Infof("setting event handles")
    networkInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: controller.enqueueNetwork,
        UpdateFunc: func(oldObj, newObj interface{}) {
            oldNetwork, newNetwork := oldObj.(simplecrdv1.Network), newObj.(simplecrdv1.Network)
            if oldNetwork.ResourceVersion == newNetwork.ResourceVersion {
                return
            }
            controller.enqueueNetwork(newObj)
        },
        DeleteFunc: controller.enqueueNetworkForDelete,
    })

    return controller
}

func (c *Controller) enqueueNetwork(obj interface{}) {
    key, err := cache.MetaNamespaceKeyFunc(obj)
    if err != nil {
        utilruntime.HandleError(err)
        return
    }
    c.workQueue.AddRateLimited(key)
}

func (c *Controller) enqueueNetworkForDelete(obj interface{}) {
    key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
    if err != nil {
        utilruntime.HandleError(err)
        return
    }
    c.workQueue.AddRateLimited(key)
}

func (c *Controller) Run(threadiness int, stopCh <- chan struct{}) error {
    defer utilruntime.HandleCrash()
    defer c.workQueue.ShutDown()

    glog.Infof("starting run control loop")

    glog.Infof("waiting to informer sync")
    if ok := cache.WaitForCacheSync(stopCh, c.netwrokSynced); ok {
        return fmt.Errorf("failed to wait informer sync")
    }

    for i := 0; i < threadiness; i++ {
        go wait.Until(c.runWorker, time.Second, stopCh)
    }

    glog.Infof("started %d worker", threadiness)

    <-stopCh
    glog.Infof("stop workers")

    return nil
}

func (c *Controller) runWorker() {
    // controll loop
    for c.ControlLoop(){
    }
}

func (c *Controller) ControlLoop() bool {
    obj, shutdown := c.workQueue.Get()
    if shutdown {
        return false
    }

    err := func(obj interface{}) error{
        defer c.workQueue.Done(obj)

        var key string
        var ok bool
        if key, ok = obj.(string); !ok {
            c.workQueue.Forget(obj)
            utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
            return nil
        }

        if err := c.syncHandler(key); err != nil {
            return fmt.Errorf("sync handler failed, %v", err)
        }

        c.workQueue.Forget(obj)
        glog.Info("Successfully handle key %s", key)
        return nil

    }(obj)

    if err != nil {
        utilruntime.HandleError(err)
        return true
    }
    return true
}



func (c *Controller) syncHandler(key string) error {
    namespace, name, err := cache.SplitMetaNamespaceKey(key)
    if err != nil {
        utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
        return nil
    }

    network, err := c.networkLister.Networks(namespace).Get(name)
    if err != nil {
        if errors.IsNotFound(err) {
            glog.Warningf("Network: %s/%s does not exist in local cache, will delete it from Neutron ...",
                namespace, name)

            glog.Infof("[Neutron] Deleting network: %s/%s ...", namespace, name)

            return nil
        } else {
            utilruntime.HandleError(fmt.Errorf("failed to list network by: %s/%s", namespace, name))
            return err
        }
    } else {
        glog.Infof("[Neutron] Try to process network: %#v ...", network)

        c.recorder.Event(network, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)

        return nil
    }
}