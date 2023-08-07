FROM golang:1.20 as builder
WORKDIR /app
COPY ./src /app
RUN go build -o main .

ADD https://github.com/benbjohnson/litestream/releases/download/v0.3.9/litestream-v0.3.9-linux-amd64-static.tar.gz litestream.tar.gz
RUN tar -xzf litestream.tar.gz -C ./

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /app/litestream /usr/local/bin/litestream
COPY litestream.yml /etc/litestream.yml
COPY start.sh /app/start.sh
RUN apk update && apk upgrade && apk add bash
CMD ["bash /app/start.sh"]

EXPOSE 8080
