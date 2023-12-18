An environment for exploring, debugging, and developing etcd.

## Usage

All you need is docker:

```
$ git clone git@github.com:jpittis/etcd-playground.git --recursive
$ cd etcd-playground
$ docker-compose up
```

- Use `docker-compose up -d --build` to rebuild etcd from source.
- You can configure latency between peers via `config/runner.yml`. By default, latency
  between instances is setup to simulate an intra-continental cluster.
- Prometheus is available on `localhost:9090`.


## Running Scenarios

The following scenario shows how you can use this playground to simulate the corruption
described in "Protocol aware recovery for consensus based storage".

Start by using docker to locate and exec to the management container:

```
$ docker ps
$ docker exec -it <container-id> /bin/bash
```

From the management container, run the following:

```
# Write a record to the cluster.
$ /etcd/bin/etcdctl --endpoints=etcd2:2379 put A 1337

# Partition etcd1 from the rest of the cluster.
$ curl -s -XPOST 'etcd1:3333/network?dev=etcd120&loss=100'
$ curl -s -XPOST 'etcd1:3333/network?dev=etcd310&loss=100'

# Write another record to the cluster.
$ /etcd/bin/etcdctl --endpoints=etcd2:2379 put B deadbeef

# Observe that etcd1 has not received the record but that etcd2 and etcd3 have.
$ curl -s 'etcd1:3333/log'
$ curl -s 'etcd2:3333/log'
$ curl -s 'etcd3:3333/log'

# Partition etcd2 from the rest of the cluster.
$ curl -s -XPOST 'etcd2:3333/network?dev=etcd120&loss=100'
$ curl -s -XPOST 'etcd2:3333/network?dev=etcd230&loss=100'

# Take down etcd3.
$ curl -s -XPOST 'etcd3:3333/etcd?enabled=false'

# Corrupt the etcd3 WAL (this needs to be run from the etcd3 container) using a perl one
# liner I ripped off of Stack Overflow.
$ cd etcd3.etcd/member/wal
$ perl -pi -e 's/deadbeef/deadbeat/g' 0000000000000000-0000000000000000.wal

# Bring etcd3 back up.
$ curl -s -XPOST 'etcd2:3333/etcd?enabled=true'
```

The result? Etcd refuses to boot... sounds like decent behavior to me:

```
etcd-playground-etcd3-1       | {"level":"fatal","ts":"2023-12-18T16:30:50.346865Z","caller":"etcdmain/etcd.go:181","msg":"discovery failed","error":"walpb: crc mismatch: expected: 5c35c817 computed: 31aa78c8: in file '0000000000000000-0000000000000000.wal' at position: 1008","stacktrace":"go.etcd.io/etcd/server/v3/etcdmain.startEtcdOrProxyV2\n\tgo.etcd.io/etcd/server/v3/etcdmain/etcd.go:181\ngo.etcd.io/etcd/server/v3/etcdmain.Main\n\tgo.etcd.io/etcd/server/v3/etcdmain/main.go:40\nmain.main\n\tgo.etcd.io/etcd/server/v3/main.go:31\nruntime.main\n\truntime/proc.go:267"}
```
