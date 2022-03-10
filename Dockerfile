FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64\
    GOPROXY="https://goproxy.cn,direct"

WORKDIR /build

COPY . .

RUN go build -o app .

WORKDIR /dist

RUN cp /build/app .

EXPOSE 8888

CMD ["/dist/app"]
