# ============ Stage 1: Build frontend ============
FROM node:22-alpine AS frontend-builder
WORKDIR /frontend/

# 使用国内镜像加速，避免网络超时
RUN npm config set registry https://registry.npmmirror.com

COPY frontend/ ./

# 安装 pnpm 并安装依赖（带重试）
RUN npm install -g pnpm@11.5.1 && \
    pnpm install --no-frozen-lockfile

RUN pnpm build

# ============ Stage 2: Build backend ============
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app/

# Go 镜像加速
RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
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
