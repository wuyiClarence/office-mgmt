FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download -x && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o platform-backend .

FROM alpine:latest

WORKDIR /app

# 复制二进制文件和配置文件到运行环境
COPY --from=builder /app/platform-backend ./

# 设置环境变量指定配置文件路径
ENV CONFIG_PATH="config"

EXPOSE 8080

CMD ["sh", "-c", "echo 'Starting app...'; ./platform-backend"]
