An environment for exploring, debugging, and developing etcd.

## Usage

All you need is docker:

```
$ git clone git@github.com:jpittis/etcd-playground.git --recursive
$ cd etcd-playground
$ docker-compose up
```

- Use `docker-compose up -d --build` to rebuild etcd from source.
- You can configure latency between peers via `config/runner.yml`.
- Client ports are available on `localhost:2379`, `localhost:3379`, and `localhost:4379`.
- Prometheus is available on `localhost:9090`.
