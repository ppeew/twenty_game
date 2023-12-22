FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /app

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/my_app /app/main.go

FROM alpine

ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/my_app /app/my_app
COPY --from=builder /app/config-pro.yaml /app/config-pro.yaml

ENTRYPOINT ["./my_app"]