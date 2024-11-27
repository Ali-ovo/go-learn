FROM golang:1.19.4-alpine3.17 AS builder

WORKDIR /build


ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod .
# 这样的好处的 docker 会调用缓存 而不会重复下载
COPY . .

RUN go env

RUN go build -ldflags="-s -w" -o hello main.go


FROM alpine

WORKDIR /build

COPY --from=builder /build/hello .

CMD ["./hello"]