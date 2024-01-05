FROM golang:latest as builder

WORKDIR /go/src/github.com/nielshojen/nanodirector/

ENV CGO_ENABLED=0 \
    GOARCH=amd64 \
    GOOS=linux

COPY . .

RUN make deps
RUN make


FROM alpine:latest

RUN apk --update add ca-certificates

COPY --from=builder /go/src/github.com/nielshojen/nanodirector/build/linux/nanodirector-linux-amd64 /usr/bin/nanodirector

EXPOSE 8000
CMD ["/usr/bin/nanodirector"]
