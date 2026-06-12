# ============ Stage 1: Build frontend ============
FROM node:22-alpine AS frontend-builder
# 限制 Node.js 堆内存为 2GB（原 4GB 过高，容易 OOM）
ENV NODE_OPTIONS="--max-old-space-size=2048"
RUN apk add --no-cache git
WORKDIR /frontend/

# 安装 pnpm 并安装依赖（若国内网络慢可使用 --network host）

COPY frontend/ ./

# 安装依赖 + 构建 + 清理：合并为一个 RUN，构建完立刻删 node_modules
# 前端产物 dist/ 是唯一需要保留的（Stage 2 通过 --from 引用）
RUN npm install -g pnpm@11.5.1 && \
    pnpm install --no-frozen-lockfile && \
    pnpm build && \
    rm -rf node_modules /root/.local/share/pnpm /root/.npm

# ============ Stage 2: Build backend ============
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app/

# GOPROXY=direct 需要 git；CGO SQLite 需要 gcc musl-dev
RUN apk add --no-cache git gcc musl-dev

ENV GOPROXY=direct \
    GOSUMDB=off \
    GONOSUMCHECK=*

# 切换 CGO SQLite 驱动，砍掉 glebarez/sqlite + modernc.org/* 整条依赖链
# 节省 ~200MB 磁盘占用

# 先复制 go.mod，再复制 Go 源码目录（不碰 frontend/）
COPY go.mod ./
COPY main.go ./
COPY cmd/ ./cmd/
COPY drivers/ ./drivers/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY public/ ./public/
COPY server/ ./server/
COPY --from=frontend-builder /frontend/dist/ ./public/dist/
RUN CGO_ENABLED=1 go build -ldflags="-w -s" -tags=jsoniter,sqlite_cgo_compat -o openlist . && \
    rm -rf /go/pkg/mod /root/.cache/go-build

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
