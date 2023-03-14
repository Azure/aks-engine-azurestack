#!/bin/bash
#
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT license.

WORKFLOW_FROM_BRANCH=${WORKFLOW_FROM_BRANCH:-master}
echo WORKFLOW_FROM_BRANCH=${WORKFLOW_FROM_BRANCH}

RELEASE_REPOSITORY=${RELEASE_REPOSITORY:-Azure/aks-engine-azurestack}
echo RELEASE_REPOSITORY=${RELEASE_REPOSITORY}

if [[ -z ${RELEASE_VERSION} ]]; then
  echo "RELEASE_VERSION is not set (e.x.: v0.76.0, v0.76.1)"
  exit 1
fi
echo RELEASE_VERSION=${RELEASE_VERSION}

if [[ -z ${RELEASE_FROM_BRANCH} ]]; then
  echo "RELEASE_FROM_BRANCH is not set (e.x: master, patch-release-v0.76.1)"
  exit 1
fi
echo RELEASE_FROM_BRANCH=${RELEASE_FROM_BRANCH}

gh workflow run create-release-branch.yaml \
  --ref ${WORKFLOW_FROM_BRANCH} \
  -R ${RELEASE_REPOSITORY} \
  -f release_version=${RELEASE_VERSION} \
  -f from_branch=${RELEASE_FROM_BRANCH}
