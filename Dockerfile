FROM golang:1.17-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 1
ENV GOOS linux
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build/go

#RUN apk add --no-cache sqlite-libs sqlite-dev && apk add --no-cache build-base
ADD go.mod .
ADD go.sum .
#RUN go mod download
COPY . .

#RUN go build -ldflags="-s -w -linkmode external -extldflags -static" -o /app/serv /build/go/cmd/server/main.go
RUN go build -ldflags="-s -w" -o /app/serv /build/go/cmd/server/main.go


FROM weblinuxgame/alpine-go:v0.1.0

WORKDIR /app

COPY --from=builder /app/serv /app/bin/queueMgrServ
# 默认配置
COPY --from=builder /build/go/etc/app.yml /app/config/app.yml


EXPOSE 80 8080
VOLUME /app/config

CMD ["/app/bin/queueMgrServ","-c","/app/config/app.yml"]
