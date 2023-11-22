We suggest cloning with `--recurse-submodules` because the dockerfile requires a local
etcd repository to build from source. Use `docker-compose up --build` to start the
playground, and `./bin/inject-latency.sh` to delay packets by 200ms each way.
