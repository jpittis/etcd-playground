#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <path> <name>"
    exit 1
fi

path="$1"
name="$2"

$path \
  --name "$name" \
  --advertise-client-urls "http://$name:2379" \
  --listen-client-urls http://0.0.0.0:2379 \
  --initial-advertise-peer-urls "http://$name:2380" \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-cluster-token etcd-cluster \
  --initial-cluster etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380 \
  --initial-cluster-state new \
  --enable-pprof \
  --logger=zap \
  --log-outputs=stderr \
  --heartbeat-interval=2000 \
  --election-timeout=10000
