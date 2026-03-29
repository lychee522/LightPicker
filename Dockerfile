# @author 肖肖雨歇 - 极致搬运版
FROM alpine:latest

# 安装基础运行库
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

# 🌟 关键：直接复制你传上去的 Linux 二进制文件
# 请确保你的文件名叫 picgo-lite-linux-amd64
COPY picgo-lite-linux-amd64 /app/picgo-lite

# 给执行权限
RUN chmod +x /app/picgo-lite

# 暴露图床端口
EXPOSE 5894

# 启动程序
CMD ["./picgo-lite"]