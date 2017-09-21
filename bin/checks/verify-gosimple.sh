#!/bin/bash

# Copyright © 2016 Samsung CNCT
#
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


# from http://github.com/kubernetes/kubernetes/hack/verify-gofmt.sh

set -o errexit
set -o nounset
set -o pipefail

ROOT=$(dirname "${BASH_SOURCE}")/../..

cd "${ROOT}"

gosimple=$(which gosimple)
if [[ ! -x "${gosimple}" ]]; then
  echo "could not find gosimple, please verify your GOPATH"
  exit 1
fi

source "${ROOT}/bin/common.sh"

errors=$( echo `packages` | xargs ${gosimple} 2>&1) || true
if [[ -n "${errors}" ]]; then
  echo "${errors}"
  exit 1
fi