# Upstream CLIProxyAPI at a release tag, with the patches in patches/ applied.
# The build stages mirror upstream's own Dockerfile.

# renovate: datasource=github-releases depName=router-for-me/CLIProxyAPI
ARG CLIPROXYAPI_VERSION=v8.0.53

FROM golang:1.26-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends build-essential git && rm -rf /var/lib/apt/lists/*

WORKDIR /app

ARG CLIPROXYAPI_VERSION
RUN git clone --depth 1 --branch "${CLIPROXYAPI_VERSION}" https://github.com/router-for-me/CLIProxyAPI.git .

COPY patches/ /patches/
RUN git apply --verbose /patches/*.patch

RUN go mod download

RUN go test ./internal/access/... ./sdk/access/...

ARG COMMIT=none
ARG BUILD_DATE=unknown

RUN CGO_ENABLED=1 GOOS=linux go build -buildvcs=false -ldflags="-s -w -X 'main.Version=${CLIPROXYAPI_VERSION}-mtaku3' -X 'main.Commit=${COMMIT}' -X 'main.BuildDate=${BUILD_DATE}'" -o ./CLIProxyAPI ./cmd/server/

FROM debian:bookworm

RUN apt-get update && apt-get install -y --no-install-recommends tzdata ca-certificates && rm -rf /var/lib/apt/lists/*

RUN mkdir /CLIProxyAPI

COPY --from=builder /app/CLIProxyAPI /CLIProxyAPI/CLIProxyAPI

COPY --from=builder /app/config.example.yaml /CLIProxyAPI/config.example.yaml

WORKDIR /CLIProxyAPI

EXPOSE 8317

ENV TZ=Asia/Shanghai

RUN cp /usr/share/zoneinfo/${TZ} /etc/localtime && echo "${TZ}" > /etc/timezone

CMD ["./CLIProxyAPI"]
