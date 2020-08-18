FROM golang:1.15.0 as builder

WORKDIR /src/
COPY . /src/

ARG RELEASE_VERSION=v0.0.1-default

ENV RELEASE_VERSION=$RELEASE_VERSION
ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -tags netgo -ldflags \
    "-w -extldflags '-static' -X main.AppVersion=${RELEASE_VERSION}" \
    -mod vendor -o ./service .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /src/service .
COPY --from=builder /src/resource ./resource/

ENTRYPOINT ["./service"]
