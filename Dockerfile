# ============ Stage 1: Build frontend ============
FROM node:22-alpine AS frontend-builder
WORKDIR /frontend/
COPY ../OpenList-Frontend/package.json ../OpenList-Frontend/pnpm-lock.yaml ../OpenList-Frontend/pnpm-workspace.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY ../OpenList-Frontend/ ./
RUN pnpm build

# ============ Stage 2: Build backend ============
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app/
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
