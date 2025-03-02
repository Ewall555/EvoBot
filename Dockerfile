FROM golang:1.23.2 AS stage-go-build

FROM node:22.11.0 AS stage-build
COPY --from=stage-go-build /usr/local/go/ /usr/local/go/
COPY --from=stage-go-build /go/ /go/
ENV GOPATH=/go
ENV PATH=/go/bin:/usr/local/go/bin:$PATH

ARG GOPROXY=https://goproxy.io
ENV CGO_ENABLED=0
ENV GO111MODULE=on

ARG TARGETARCH
ARG NPM_REGISTRY="https://registry.npmmirror.com"
ENV NPM_REGISTY=$NPM_REGISTRY

RUN set -ex \
    && npm config set registry ${NPM_REGISTRY} \
    && yarn config set registry ${NPM_REGISTRY}

WORKDIR /opt/evobot

COPY . .

RUN cd frontend && yarn install && yarn build

WORKDIR /opt/evobot/cmd/server

ARG GOPROXY=https://goproxy.io
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux

RUN --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download -x \
    && go build -o /opt/evobot/evobot ./main.go

FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /opt/evobot

COPY --from=stage-build /opt/evobot/evobot .

EXPOSE 24916

CMD ["/opt/evobot/evobot"]