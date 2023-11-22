#!/bin/bash

docker exec etcd-playground-etcd1-1 tc qdisc del dev eth0 root
docker exec etcd-playground-etcd2-1 tc qdisc del dev eth0 root
docker exec etcd-playground-etcd3-1 tc qdisc del dev eth0 root
