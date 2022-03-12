# Simple Controller Custom Resource

This is a template of custom APIs controller using code-generator


## Usage

Run go get code-generator
```shell script
gp get  k8s.io/code-generator@kubenetes-$KUBENETES_VERSION
``` 
Modify follows: 
1. Modify the directory name in simplecrd
1. Modify Group name and version in [register.go](pkg/apis/simplecrd/v1)
1. Modify struct name and spec [type.go](pkg/apis/simplecrd/v1/type.go)

Fioally, run this script
```shell script
./hack/update-codegen.sh
```


