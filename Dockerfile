FROM golang:1.15 as builder
WORKDIR /go/src/github.com/irotoris/go-ipadder-server/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ipaddr-server

FROM alpine:3.12
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/irotoris/go-ipadder-server/ipaddr-server /ipaddr-server

ENTRYPOINT [ "/ipaddr-server" ]
EXPOSE 8080
