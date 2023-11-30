# 使用 Go 镜像作为基础镜像
FROM golang:1.21.4 AS builder

# 设置工作目录
WORKDIR /app

# 将代码复制到容器中
COPY . .

# 构建可执行文件
RUN go build -o myapp

# 创建最终的小型镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 复制可执行文件从builder镜像
COPY --from=builder /app/myapp .

# 暴露应用程序端口
EXPOSE 8888

# 运行应用程序
CMD ["./myapp"]