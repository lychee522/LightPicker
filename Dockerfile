# @author 肖肖雨歇
# 第一阶段：编译前端 Vue
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ .
RUN npm run build

# 第二阶段：编译 Go 后端
FROM golang:1.22-alpine AS backend-builder
# 禁用 CGO，确保生成纯净的静态二进制文件，跨平台通杀！
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 把第一阶段编译好的前端 dist 复制过来供 go:embed 吞噬
COPY --from=frontend-builder /app/web/dist /app/web/dist
# 暴力压缩编译参数 (-s -w 去掉调试信息，体积减半)
RUN go build -ldflags="-s -w" -o picgo-lite cmd/main.go

# 第三阶段：最终极简运行环境 (1C1G 小鸡狂喜)
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
# 只拷贝那个编译出来的终极二进制文件！
COPY --from=backend-builder /app/picgo-lite /app/picgo-lite
# 声明存储目录，防止数据随容器消亡
VOLUME ["/app/storage"]
EXPOSE 8080
CMD ["./picgo-lite"]