# Build
FROM golang:latest AS build
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /work

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o /app .

# Deploy
FROM alpine:latest
WORKDIR /work
COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=build /app ./
EXPOSE 8080
ENTRYPOINT ["./app"]
