# ============ Stage 1: Build frontend ============
FROM node:22-alpine AS frontend-builder
ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN apk add --no-cache git
WORKDIR /frontend/

# 安装 pnpm 并安装依赖（若国内网络慢可使用 --network host）

COPY frontend/ ./

# 安装 pnpm 并安装依赖（带重试）
RUN npm install -g pnpm@11.5.1 && \
    pnpm install --no-frozen-lockfile

RUN pnpm build

# ============ Stage 2: Build backend ============
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app/

ENV GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=off \
    GONOSUMCHECK=*

RUN apk add --no-cache gcc musl-dev
COPY ./ ./

# 清理已删除驱动的残留依赖，然后下载
RUN go mod tidy && go mod download
COPY --from=frontend-builder /frontend/dist/ ./public/dist/
RUN CGO_ENABLED=1 go build -ldflags="-w -s" -tags=jsoniter -o openlist .

# ============ Stage 3: Runtime ============
FROM alpine:edge
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /opt/openlist/
COPY --from=backend-builder /app/openlist ./
COPY --chmod=755 entrypoint.sh /entrypoint.sh

ENV UMASK=022 TZ=Asia/Shanghai
VOLUME /opt/openlist/data/
EXPOSE 5244
CMD [ "/entrypoint.sh" ]
