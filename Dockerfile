FROM golang:latest as builder
WORKDIR /app
COPY ./src /app
RUN go build -o main

FROM alpine:latest as production
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/main /app/main
CMD ["./main"]
