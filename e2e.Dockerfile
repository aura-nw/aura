FROM golang:1.18.4-alpine3.15 as builder

RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git

COPY . /aura
WORKDIR /aura

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.1/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.1/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 9ecb037336bd56076573dc18c26631a9d2099a7f2b40dc04b6cae31ffb4c8f9a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 6e4de7ba9bad4ae9679c7f9ecf7e283dd0160e71567c6a7be6ae47c81ebe7f32

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp "/lib/libwasmvm_muslc.$(uname -m).a" /lib/libwasmvm_muslc.a

RUN make clean && make build VERBOSE=1 BUILD_TAGS=muslc

RUN echo "Ensuring binary is statically linked ..." \
  && (file /aura/build/aurad | grep "statically linked")

FROM golang:1.18.4-bullseye as ignite
ARG BUILD_ENV=dev

RUN curl https://get.ignite.com/cli@v0.24.0! | bash

COPY . /aura
WORKDIR /aura

COPY config.yml config.yml
RUN ignite chain init

FROM alpine:3.15

COPY --from=ignite /root/.aura /root/.aura
COPY --from=builder /aura/build/aurad /usr/bin/aurad

# rest grpc p2p rpc
EXPOSE 1317 9090 26656 26657

CMD ["/usr/bin/aurad", "start"]