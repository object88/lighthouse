#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

cd $(dirname ${BASH_SOURCE})

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
# CODEGEN_PKG=${CODEGEN_PKG:-./vendor/k8s.io/code-generator}

# TEMPDIR=$(mktemp -d)

# bash vendor/k8s.io/code-generator/generate-groups.sh all \
#   github.com/object88/lighthouse/pkg/k8s/client \
#   github.com/object88/lighthouse/pkg/k8s/apis \
#   engineering.lighthouse:v1alpha1 \
#   --go-header-file "${SCRIPT_ROOT}/buildtools/custom-boilerplate.go.txt" \
#   --output-base "${TEMPDIR}"

CRD_OPTIONS="crd:trivialVersions=true,preserveUnknownFields=false"

../bin/controller-gen ${CRD_OPTIONS} rbac:roleName=manager-role webhook paths="../..." output:crd:artifacts:config="../config/crd/bases"

../bin/controller-gen object:headerFile="$PWD/boilerplate.go.txt" paths="../..."


# echo "generation complete"

# if [ -n "$(ls -A "${SCRIPT_ROOT}/pkg/k8s/apis")" ]; then
#   echo "copying apis"
#   rsync -a "${TEMPDIR}/github.com/object88/lighthouse/pkg/k8s/apis/" "${SCRIPT_ROOT}/pkg/k8s/apis"
# fi
# if [ -n "$(ls -A "${SCRIPT_ROOT}/pkg/k8s/client")" ]; then
#   echo "copying client"
#   rsync -a "${TEMPDIR}/github.com/object88/lighthouse/pkg/k8s/client/" "${SCRIPT_ROOT}/pkg/k8s/client"
# fi

# echo "done"
