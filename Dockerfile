FROM golang:1.22-alpine as builder

RUN set -xe  && \
        apk update && apk upgrade  && \
        apk add --no-cache make git
ENV CGO_ENABLED=0
COPY . /go/src/app
WORKDIR /go/src/app

RUN make build && \
    cp ./build/server /server && \
    cp ./build/client /client

FROM alpine
COPY --from=builder /server /server
COPY --from=builder /client /client
COPY --from=builder /go/src/app/assets /assets
CMD ["/server"]
