#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-bulletin-api
  make audit
popd
