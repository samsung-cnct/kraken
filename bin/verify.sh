#!/bin/bash

# Copyright © 2016 Samsung CNCT
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


# using github.com/kubernetes/kubernetes/hack/make-rules/verify.sh as basis for this file.

set -o errexit
set -o nounset
set -o pipefail

ROOT=$(dirname "${BASH_SOURCE}")/..
source "${ROOT}/bin/common.sh"


# Collect Failed tests in this Array , initialize to nil
FAILED_TESTS=()

function print-failed-tests {
  echo -e "========================"
  echo -e "${color_red}FAILED TESTS${color_norm}"
  echo -e "========================"
  for t in ${FAILED_TESTS[@]}; do
      echo -e "${color_red}${t}${color_norm}"
  done
}

function run-cmd {
  if ${SILENT}; then
    "$@" &> /dev/null
  else
    "$@"
  fi
}

function run-check {
  local -r check=$1
  local -r runner=$2

  echo -e "Verifying ${check}"
  local start=$(date +%s)
  run-cmd "${runner}" "${check}" && tr=$? || tr=$?
  local elapsed=$(($(date +%s) - ${start}))
  if [[ ${tr} -eq 0 ]]; then
    echo -e "${color_green}SUCCESS${color_norm}  ${check}\t${elapsed}s"
  else
    echo -e "${color_red}FAILED${color_norm}   ${check}\t${elapsed}s"
    ret=1
    FAILED_TESTS+=(${check})
  fi
}

SILENT=true

while getopts ":v" opt; do
  case ${opt} in
    v)
      SILENT=false
      ;;
    \?)
      echo "Invalid flag: -${OPTARG}" >&2
      exit 1
      ;;
  esac
done

if ${SILENT} ; then
  echo "Running in silent mode, run with -v if you want to see script logs."
fi


ret=0
run-check "${ROOT}/bin/checks/verify-go-vet.sh" bash
run-check "${ROOT}/bin/checks/verify-gofmt.sh" bash
run-check "${ROOT}/bin/checks/verify-golint.sh" bash
run-check "${ROOT}/bin/checks/verify-gosimple.sh" bash


if [[ ${ret} -eq 1 ]]; then
    print-failed-tests
fi
exit ${ret}