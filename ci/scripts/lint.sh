#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-bulletin-api
  make lint
popd
