#!/bin/bash -eux

pushd dp-bulletin-api
  make build
  cp build/dp-bulletin-api Dockerfile.concourse ../build
popd
