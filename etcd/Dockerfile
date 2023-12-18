FROM alpine:latest as builder
WORKDIR /etcd

# We require make, bash, and Go to build etcd. In the future, we may want to switch from
# alpine over to the standard Go image because downloading Go takes forever.
RUN apk add --no-cache make bash go

# We start by only copying over the go.mod and go.sum files so we can trigger a download
# of all the project's dependencies, and cache it in a layer.
COPY ./etcd/etcd/go.* .
COPY ./etcd/etcd/api/go.* ./api/
COPY ./etcd/etcd/server/go.* ./server/
COPY ./etcd/etcd/pkg/go.* ./pkg/
COPY ./etcd/etcd/client/pkg/go.* ./client/pkg/
COPY ./etcd/etcd/client/internal/v2/go.* ./client/internal/v2/
COPY ./etcd/etcd/client/v3/go.* ./client/v3/
COPY ./etcd/etcd/etcdctl/go.* ./etcdctl/
COPY ./etcd/etcd/etcdutl/go.* ./etcdutl/
COPY ./etcd/etcd/tests/go.* ./tests/
COPY ./etcd/etcd/tools/mod/go.* ./tools/mod/
RUN go mod download

# Now we build both the standard etcd binaries (etcd, etcdctl, etcdutl), as well as the
# tools (benchmark, etcd-dump-db, etcd-dump-logs, local-tester). The mount is for reusing
# Go's object cache.
ADD ./etcd/etcd .
RUN --mount=type=cache,target=/root/.cache/go-build make build
RUN --mount=type=cache,target=/root/.cache/go-build make tools

FROM golang:latest as runner-builder
WORKDIR /runner

# Copy over go.sum to cache runner dependencies.
COPY ./etcd/runner/go.* .
RUN go mod download

ADD ./etcd/runner .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build

# Finally, we copy all the binaries over to a fresh container.
FROM alpine:latest
WORKDIR /etcd

# We require iproute2 for iptables and tc.
RUN apk add --no-cache iproute2 bash curl perl

COPY --from=builder ./etcd/bin ./bin/
COPY --from=runner-builder ./runner/runner ./bin/runner
