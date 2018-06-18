#!/bin/bash

CONTAINER_ID="$(kubectl get pods | grep osqueryi | cut -d " " -f1)"
if [ -z "$CONTAINER_ID" ]; then
    echo "No container found"
    exit 1
fi
kubectl exec -it $CONTAINER_ID -- osqueryi --extension ./osquery-kubelet
