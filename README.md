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

## 快速开始

### 直接运行

```bash
# 1. 编译
go build -ldflags="-w -s" -tags=jsoniter -o openlist .

# 2. 准备数据目录
mkdir -p data

# 3. 启动
./openlist server --debug
```

### Docker 部署

**推荐方式 (docker compose)**，`host` 网络模式可避免端口映射问题：

```yaml
# docker-compose.yml
services:
  openlist:
    restart: always
    network_mode: host
    volumes:
      - './data:/opt/openlist/data'
    environment:
      - TZ=Asia/Shanghai
    image: openlistteam/openlist:latest
```

```bash
docker compose up -d
```

**bridge 模式**：

```yaml
services:
  openlist:
    restart: always
    ports:
      - '5244:5244'
    volumes:
      - './data:/opt/openlist/data'
    environment:
      - TZ=Asia/Shanghai
    image: openlistteam/openlist:latest
```

> 注意：bridge 模式下访问 `localhost` 时实际指向容器内部。如需挂载宿主机本地目录，应使用 `network_mode: host` 或挂载对应卷。

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
