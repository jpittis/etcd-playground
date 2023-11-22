#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <delay>"
    exit 1
fi

delay="$1"

docker exec etcd-playground-etcd1-1 tc qdisc add dev eth0 root netem delay "$delay"
docker exec etcd-playground-etcd2-1 tc qdisc add dev eth0 root netem delay "$delay"
docker exec etcd-playground-etcd3-1 tc qdisc add dev eth0 root netem delay "$delay"
