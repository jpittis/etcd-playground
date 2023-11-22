#!/bin/bash

docker exec etcd-playground-etcd1-1 tc qdisc add dev eth0 root netem delay 200ms
docker exec etcd-playground-etcd2-1 tc qdisc add dev eth0 root netem delay 200ms
docker exec etcd-playground-etcd3-1 tc qdisc add dev eth0 root netem delay 200ms
