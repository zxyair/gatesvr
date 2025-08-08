FROM ubuntu:latest
LABEL authors="xyxvzhang"

ENTRYPOINT ["top", "-b"]

# 使用官方Go镜像作为构建环境
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制依赖文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目文件
COPY . .

# 构建可执行文件（禁用CGO，减小镜像大小）
RUN CGO_ENABLED=0 GOOS=linux go build -o client

# 使用轻量级运行时镜像
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制可执行文件
COPY --from=builder /app/server .

# 暴露服务端口（根据实际端口修改）
EXPOSE 3353

# 设置容器启动命令
CMD ["./client"]
