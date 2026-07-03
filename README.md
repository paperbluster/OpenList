# OpenList

局域网多存储文件管理程序，将本地目录、WebDAV、远端 OpenList 统一挂载到一个目录树下，通过 Web 界面管理。

## 用途

- 局域网内集中管理多个 NAS / 文件服务器
- 将分散的 WebDAV、OpenList 远端聚合为统一入口
- 轻量部署，适合家庭/小型办公环境

## 支持的存储

| 类型 | 驱动 |
|---|---|
| 本地 | 本地目录 |
| 协议 | WebDAV |
| 远端 | OpenList, OpenList 分享, AList v3 |
| 虚拟 | 别名(快捷方式), 虚拟聚合, 自动索引 |

## 项目结构

```
├── cmd/            ← Go 命令行
├── drivers/        ← 存储驱动
├── internal/       ← 核心逻辑
├── server/         ← HTTP 服务
├── frontend/       ← 前端源码 (Vite + Solid.js)，构建产物 → public/dist/
├── public/         ← 前端静态资源 + 内嵌资源
└── main.go         ← 入口
```

## 快速开始

```bash
# 1. 构建前端
cd frontend
pnpm install && pnpm build

# 2. 编译并启动后端
cd ..
go build -ldflags="-w -s" -tags=jsoniter -o openlist .
mkdir -p data
./openlist server
```

访问：**http://localhost:5244**，管理后台：**http://localhost:5244/@manage**

## Docker 部署

```bash
# 编译
docker build --network host -t openlist .

# 启动（host 网络模式，推荐局域网使用）
docker run -d \
  --name openlist \
  --network host \
  --restart always \
  -v /mnt/data:/opt/openlist/data \
  -e TZ=Asia/Shanghai \
  openlist

# 或桥接模式
docker run -d \
  --name openlist \
  -p 5244:5244 \
  --restart always \
  -v /mnt/alist:/opt/openlist/data \
  -e TZ=Asia/Shanghai \
  openlist
```

常用管理命令：
```bash
docker logs -f openlist    # 查看日志
docker restart openlist    # 重启
docker stop openlist       # 停止
docker exec -it openlist sh  # 进入容器
```

## 挂载存储

1. 访问管理后台，点击 `存储` → `添加`
2. 选择驱动类型，填写挂载路径和参数
3. 保存后即可在主页浏览

常见配置：

| 存储类型 | 关键参数 |
|---|---|
| 本地目录 | 挂载路径 + 本地文件夹路径 |
| WebDAV | 挂载路径 + WebDAV 地址 + 用户名密码 |
| OpenList 远端 | 挂载路径 + 远端地址 + 用户名密码 |
| AList v3 远端 | 挂载路径 + 远端地址 + 用户名密码 |

支持运行时动态添加/修改/删除挂载。

## 命令行

```
./openlist server      # 前台启动
./openlist admin       # 管理管理员账号
./openlist --help      # 查看所有命令
```
