FROM golang:1.24.5-alpine AS builder

WORKDIR /demo

COPY go.mod go.sum ./

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /demo/myapp .

EXPOSE 1234

CMD ["./myapp"]