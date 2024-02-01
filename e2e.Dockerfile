FROM golang:1.21-alpine3.17 AS go-builder
#ARG arch=x86_64

RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git

COPY . /aura
WORKDIR /aura

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.4.2/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.4.2/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
#RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 9ecb037336bd56076573dc18c26631a9d2099a7f2b40dc04b6cae31ffb4c8f9a
#RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 6e4de7ba9bad4ae9679c7f9ecf7e283dd0160e71567c6a7be6ae47c81ebe7f32

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp "/lib/libwasmvm_muslc.$(uname -m).a" /lib/libwasmvm_muslc.a

RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build

RUN echo "Ensuring binary is statically linked ..." \
  && (file /aura/build/aurad | grep "statically linked")

FROM golang:1.21-bullseye as ignite

RUN curl https://get.ignite.com/cli@v0.27.1! | bash

COPY . /aura
WORKDIR /aura

COPY config.yml config.yml
RUN ignite chain init

FROM alpine:3.17

COPY --from=ignite /root/.aura /root/.aura
COPY --from=go-builder /aura/build/aurad /usr/bin/aurad

# rest grpc p2p rpc
EXPOSE 1317 9090 26656 26657

CMD ["/usr/bin/aurad", "start"]
