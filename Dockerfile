# 第一阶段：构建前端
FROM node:22.11.0 AS frontend-build

ARG NPM_REGISTRY="https://registry.npmmirror.com"
ENV NPM_REGISTRY=${NPM_REGISTRY}

WORKDIR /app

COPY frontend ./frontend

RUN npm config set registry ${NPM_REGISTRY} \
    && yarn config set registry ${NPM_REGISTRY} \
    && cd frontend \
    && yarn install \
    && yarn build


# 第二阶段：构建 Go
FROM golang:1.23.2 AS go-build

ARG GOPROXY=https://goproxy.io
ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOPROXY=${GOPROXY} 

WORKDIR /opt/evobot

# 先复制所有源码
COPY . .

# 再复制前端构建产物，让 embed 可以找到
COPY --from=frontend-build /app/frontend/dist ./frontend/dist

WORKDIR /opt/evobot/cmd

RUN go mod download -x && go build -ldflags="-s -w" -o /opt/evobot/evobot ./main.go


# 第三阶段：构建最终产物
FROM alpine:latest AS final

RUN apk add --no-cache bash tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /opt/evobot

# 拷贝可执行文件与配置
COPY --from=go-build /opt/evobot/evobot .
COPY --from=go-build /opt/evobot/cmd/conf/config-example.yaml ./config.yaml

# 可选：再拷贝一次前端资源（如果运行时也要 serve 静态文件）
COPY --from=frontend-build /app/frontend/dist ./frontend/dist

ENV TZ=Asia/Shanghai

EXPOSE 24916

CMD ["/opt/evobot/evobot", "-c", "/opt/evobot/conf/config.yaml"]
