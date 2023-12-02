An environment for exploring, debugging, and developing etcd.

## Usage

- All you need is docker:

```
$ git clone git@github.com:jpittis/etcd-playground.git --recursive
$ cd etcd-playground
$ docker-compose up
```

- Use `docker-compose up -d --build` to rebuild etcd from source.
- You can also inject/clear latency between peers:

```
$ ./bin/inject-latency.sh
$ time etcdctl get example
etcdctl get example  0.02s user 0.01s system 1% cpu 1.444 total
$ ./bin/clear-latency.sh
$ time etcdctl get example
etcdctl get example  0.02s user 0.01s system 93% cpu 0.031 total
```

- Client ports are available on `localhost:2379`, `localhost:3379`, and `localhost:4379`.
- Prometheus is available on `localhost:9090`.
