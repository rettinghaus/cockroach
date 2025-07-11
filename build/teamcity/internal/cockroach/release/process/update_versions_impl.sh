#!/usr/bin/env bash

# Copyright 2022 The Cockroach Authors.
#
# Use of this software is governed by the CockroachDB Software License
# included in the /LICENSE file.


set -xeuo pipefail

to=dev-inf+release-dev@cockroachlabs.com
dry_run=true
version_bump_only=false
# override dev defaults with production values
if [[ -z "${DRY_RUN}" ]] ; then
  echo "Setting production values"
  to=release-engineering-team@cockroachlabs.com
  dry_run=false
fi

if [[ -n "${VERSION_BUMP_ONLY}" ]] ; then
  version_bump_only=true
fi

# run git fetch in order to get all remote branches
git fetch --tags -q origin

# install gh and helm
curl -fsSL -o /tmp/gh.tar.gz https://github.com/cli/cli/releases/download/v2.32.1/gh_2.32.1_linux_amd64.tar.gz
echo "5c9a70b6411cc9774f5f4e68f9227d5d55ca0bfbd00dfc6353081c9b705c8939  /tmp/gh.tar.gz" | sha256sum -c -
tar --strip-components 1 -xf /tmp/gh.tar.gz
curl -fsSL -o /tmp/helm.tar.gz https://get.helm.sh/helm-v3.14.1-linux-amd64.tar.gz
echo "75496ea824f92305ff7d28af37f4af57536bf5138399c824dff997b9d239dd42  /tmp/helm.tar.gz" | sha256sum -c -
tar -C bin --strip-components 1 -xf /tmp/helm.tar.gz linux-amd64/helm
export PATH=$PWD/bin:$PATH

bazel build --config=crosslinux //pkg/cmd/release

$(bazel info --config=crosslinux bazel-bin)/pkg/cmd/release/release_/release \
  update-versions \
  --dry-run=$dry_run \
  --version-bump-only=$version_bump_only \
  --released-version=$RELEASED_VERSION \
  --next-version=$NEXT_VERSION \
  --template-dir=pkg/cmd/release/templates \
  --smtp-user=cronjob@cockroachlabs.com \
  --smtp-host=smtp.gmail.com \
  --smtp-port=587 \
  --artifacts-dir=/artifacts \
  --to=$to
