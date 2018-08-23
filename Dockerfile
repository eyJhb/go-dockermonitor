FROM golang:1.10.3 as builder
RUN go get -v github.com/eyJhb/go-dockermonitor
WORKDIR /go/src/github.com/eyJhb/go-dockermonitor
RUN GOOS=linux CGO_ENABLED=0 go build -o app .

FROM alpine:3.7
ENV VERSION 1.0
LABEL maintainer="eyjhb <eyjhbb@gmail.com>"
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/eyJhb/go-dockermonitor/app .


ENTRYPOINT ["./app"]

