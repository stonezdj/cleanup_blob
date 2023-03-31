FROM golang:1.19.7-alpine3.17 AS builder

COPY main.go /src/main.go
COPY go.mod /src/go.mod
COPY go.sum /src/go.sum


RUN  cd /src && go mod tidy && go build -o /bin/cleanup_blob main.go

FROM alpine:3.17

COPY --from=builder /bin/cleanup_blob /bin/cleanup_blob

CMD ["/bin/cleanup_blob"]
