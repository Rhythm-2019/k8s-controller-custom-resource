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
package v1

import (
    "github.com/Rhythm-2019/k8s-controller-custom-resource/pkg/apis/simplecrd"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = &schema.GroupVersion{
    Group:   simplecrd.GroupNmae,
    Version: simplecrd.Version,
}

var (
    SchemaBuilder = runtime.NewSchemeBuilder(addKnownTypes)
    AddToSchema = SchemaBuilder.AddToScheme
)

func Resource(resource string) schema.GroupResource {
    return SchemeGroupVersion.WithResource(resource).GroupResource()
}
func Kind(kind string) schema.GroupKind {
    return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func addKnownTypes(scheme *runtime.Scheme) error {

    scheme.AddKnownTypes(
        *SchemeGroupVersion,
        Network{},
        NetworkList{},
    )

    metav1.AddToGroupVersion(scheme, *SchemeGroupVersion)
    return nil


}