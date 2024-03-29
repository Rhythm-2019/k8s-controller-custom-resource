#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

# modify follow env var
KUBENETES_VERSION="1.17.5"
BASE_PACKAGE="github.com/Rhythm-2019/k8s-controller-custom-resource"
GENERATORS="all"
CUSTOM_RESOURCE_NAME="simplecrd"
CUSTOM_RESOURCE_VERSION="v1"

cd `dirname $0`
chmod +x ../vendor/k8s.io/code-generator/generate-groups.sh

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
../vendor/k8s.io/code-generator/generate-groups.sh \
  $GENERATORS \
  $BASE_PACKAGE/pkg/client \
  $BASE_PACKAGE/pkg/apis \
  $CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION \
  --go-header-file $(pwd)/boilerplate.go.txt \
