# Sample Controller Custom Resource

This is a template of custom APIs controller with code-generator. 

## Usage

Run go get code-generator
```shell script
$ go get k8s.io/code-generator@kubenetes=$KUBENETES_VERSION
``` 
Modify follows: 
1. Modify the directory name in simplecrd
1. Modify Group name and version in [register.go](pkg/apis/samplecrd/v1)
1. Modify struct name and spec [type.go](pkg/apis/samplecrd/v1/type.go)
4. Modify [update-codegen.sh](/hack/update-codegen.sh) KUBENETES_VERSION and BASE_PACKAGE, run it and copy to pkg/client,

Buikd code and image: 
1. Build execute file by ```go build```
2. Build image by ```docker build -t <image-name>:<tag-name> .```, push to hub
3. Create [Role, RoleBinding, ServiceAccount](/k8s/rbac.yaml), [CRD](/k8s/crd.yaml) and [Pod](/k8s/controller.yaml)

See the effect:
Create or delete [Example-Resource](example/example-network.yaml), see log
```text
W0825 13:31:40.924296       1 controller.go:182] Network: default/example-network does not exist in local cache, will delete it from Neutron ...
I0825 13:31:40.925713       1 controller.go:185] [Neutron] Deleting network: default/example-network ...
I0825 13:31:40.925726       1 controller.go:157] Successfully handle key default/example-network
``` 
 


