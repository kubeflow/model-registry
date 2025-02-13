#!/bin/sh
set -o errexit

# Some of this copy-pasted from https://kind.sigs.k8s.io/docs/user/local-registry/

# 1. Create registry container unless it already exists
reg_name='local-mr-registry'
reg_port='5001'
if [ "$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --network bridge --name "${reg_name}" \
    registry:2
fi
