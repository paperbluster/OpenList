# OpenList

多存储文件列表程序，支持将各类网盘和存储协议统一挂载到一个目录树下，通过 Web 界面管理。

## 支持的存储

| 类型 | 驱动 |
|---|---|
| 标准协议 | FTP, SFTP, SMB/CIFS, WebDAV |
| 对象存储 | S3 (及兼容), Azure Blob, 又拍云 USS |
| 开放平台 | OpenList, AList v3, Cloudreve v3/v4 |
| 网盘 | PikPak, PikPak 分享 |
| 本地/虚拟 | 本地存储, 别名(快捷方式), 加密存储, URL 树, .strm 流媒体索引, 虚拟聚合 |

## 项目结构

```
个人自研项目/
├── OpenList/          ← 后端 (Go)
└── OpenList-Frontend/ ← 前端 (Vite + Solid.js)
```

前后端是两个独立项目。后端编译后内嵌前端编译产物（`public/dist/`）。

## 快速开始（本地开发）

```bash
# 1. 构建前端
cd ../OpenList-Frontend
pnpm install && pnpm build

# 2. 把前端产物复制到后端
cp -r dist/ ../OpenList/public/dist/

# 3. 编译并启动后端
cd ../OpenList
go build -ldflags="-w -s" -tags=jsoniter -o openlist .
mkdir -p data
./openlist server --debug
```

> 访问：**http://localhost:5244**，管理后台：**http://localhost:5244/@manage**

## Docker 部署

Dockerfile 会自动编译前端+后端，一步打包。

**项目结构需要保持：** 两个目录在同一个父目录下，因为 Dockerfile 的 context 需要同时访问前后端源码。

```
some-dir/
├── OpenList/             ← docker compose 在这里执行
│   ├── Dockerfile
│   └── docker-compose.yml
└── OpenList-Frontend/    ← 前端源码
```

### host 模式（推荐）

```bash
cd OpenList
docker compose up -d
```

docker-compose.yml 已配置 `network_mode: host`，直接访问宿主机端口 `5244`。

### bridge 模式

```yaml
services:
  openlist:
    build:
      context: ..
      dockerfile: OpenList/Dockerfile
    restart: always
    ports:
      - '5244:5244'
    volumes:
      - './data:/opt/openlist/data'
    environment:
      - TZ=Asia/Shanghai
```

> bridge 模式下如需挂载宿主机本地文件，需要额外 `-v /host/path:/container/path` 映射。

## 访问

启动后访问：**`http://localhost:5244`**

首次使用需设置管理员密码，访问管理后台：**`http://localhost:5244/@manage`**

## 挂载存储

1. 访问管理后台，点击左侧 `存储` → `添加`
2. 选择驱动类型，填写必需参数（挂载路径必填，其他根据驱动提示填写）
3. 保存后即可在主页浏览该存储

常见配置示例：

| 存储类型 | 关键参数 |
|---|---|
| 本地目录 | 挂载路径 + 本地文件夹路径 |
| S3 | `AccessKey`, `SecretKey`, `Endpoint`, `Bucket` |
| WebDAV | 挂载路径 + WebDAV 地址 + 用户名密码 |
| SFTP | 挂载路径 + 主机地址 + 端口 + 用户名密码 |
| FTP | 挂载路径 + 主机地址 + 端口 + 用户名密码 |
| SMB | 挂载路径 + 共享地址 + 用户名密码 |
| Azure Blob | 挂载路径 + AccessKey + 容器名 |
| PiKpak | 挂载路径 + 用户名密码 |

所有参数可通过 API 或管理界面配置，支持在运行中动态添加/修改/删除挂载。

## 命令行

```
./openlist server          # 前台启动
./openlist start           # 后台守护启动
./openlist stop            # 停止
./openlist restart         # 重启
./openlist admin           # 管理管理员账号
./openlist --help          # 查看所有命令
```
