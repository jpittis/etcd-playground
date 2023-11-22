etcd-playground is an environment for exploring, debugging, and developing etcd.

## Usage

```
$ git clone git@github.com:jpittis/etcd-playground.git --recursive
$ cd etcd-playground
$ docker-compose up
$ ./bin/inject-latency.sh
```

Use `docker-compose up -d --build` to rebuild etcd from source.
