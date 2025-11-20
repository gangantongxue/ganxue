FROM golang:1.24.1 AS builder

WORKDIR /app
COPY . .

ENV GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

FROM alpine:latest

RUN apk update
RUN apk add --no-cache \
    go \
    musl-dev \
    netcat-openbsd

WORKDIR /home/ganxue/ganxue-server

COPY --from=builder /app/ .
COPY --from=builder /app/server .
RUN ls -al
RUN chmod +x wait_for_services.sh
RUN chmod +x server

CMD ["./server"]