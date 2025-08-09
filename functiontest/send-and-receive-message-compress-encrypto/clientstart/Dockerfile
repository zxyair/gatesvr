# 构建阶段
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制 gatesvr 下所有内容（假设你在 gatesvr 根目录执行 docker build）
COPY . .

# 编译 clientstart 目录下的 main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o client ./functiontest/send-and-receive-message-compress-encrypto/clientstart

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 复制可执行文件
COPY --from=builder /app/client .

# 复制配置文件（如果需要）
COPY --from=builder /app/functiontest/send-and-receive-message-compress-encrypto/clientstart/ .


CMD ["./client"]