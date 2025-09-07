# StealthIM DBGateway

数据库网关 `0.0.1`

Redis / Mysql 网关

> `.proto` 文件：`./proto/db_gateway.proto`

## 构建

### 依赖

Go 版本：`1.24.0`

软件包：`protobuf` `protobuf-dev` `make`

> 命令行工具 `protoc` `make`(gnumake)

```bash
go mod download
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 命令行

```bash
make # 构建并运行
make build # 构建可执行文件

# 构建指定环境
make build_windows
make build_linux
make build_docker

make release # 构建所有平台

make proto # 生成 proto

make clean # 清理
```

## 配置

默认会读取当前文件夹 `config.toml` 文件（不存在会自动生成模板）

> docker 镜像为 `cfg/config.toml`

默认配置：

```toml
[mysql]
host = "127.0.0.1"           # Mysql地址
port = 3306                  # Mysql端口
maxconn = 50                 # 最大连接数
minconn = 10                 # 最小连接数
user = "root"                # 用户名，首次建议使用root，之后可以创建新用户
password = "<YOUR_PASSWORD>" # 密码
prefix = ""                  # 数据库前缀

#     下面单独设置每一个数据库的配置
#   若某项设置未填写则使用上方全局配置
# [mysql.db_msg]
# [mysql.db_users]
# [mysql.masterdb]
# [mysql.logging]
# [mysql.groups]
# [mysql.file]
# [mysql.session]
# host = "127.0.0.1"
# port = 3306
# maxconn = 50
# minconn = 10
# user = "root"
# password = "<YOUR_PASSWORD>"

[redis]
host = "127.0.0.1"           # Redis地址
port = 6379                  # Redis端口
dbname = 0                   # 数据库编号，通常无需改动
password = "<YOUR_PASSWORD>" # 密码
cachetime = 3600             # 缓存时间，单位秒

[grpc]
host = "127.0.0.1" # GRPC地址
port = 50051       # GRPC监听端口
log = false        # 启用日志，调试功能，上线建议关闭
```

也可使用 `--config={PATH}` 参数指定配置文件路径
