### Quick Installation
```shell
go install github.com/1cool/w@latest
```
### Usage
1、run init command
```shell
wo init your_project

init successful
please cd your project dirname,then edit your config file `config.yaml`
finally please run `go mod tidy && go run main.go`
```

2、cd your_project. then run command `go mod tidy`

3、config config.yaml
```shell
use mysql command 
create database your_project;

vim config.yaml

database:
    driver: mysql # 数据库驱动；支持mysql
    mysql:
        database: your_project # 数据库名称
        host: 127.0.0.1 # 数据库地址
        port: 3306 # 端口
        username: root # 账号
        password: # 密码
http:
    addr: 127.0.0.1:8000 # 启动地址
```
4、run new entity command for `crud business`. should generate model、ent schema、repository、service、handler、request、response.
```shell
wo new entity user

new entity successful user
```
5、start project
```shell
go run main.go

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/v1/example           --> your_project/internal/httptransport.(*handler).example-fm (1 handlers)
[GIN-debug] POST   /api/v1/users             --> your_project/internal/httptransport.(*handler).AddUser-fm (1 handlers)
[GIN-debug] PUT    /api/v1/users/:id         --> your_project/internal/httptransport.(*handler).UpdateUser-fm (1 handlers)
[GIN-debug] GET    /api/v1/users/:id         --> your_project/internal/httptransport.(*handler).ShowUser-fm (1 handlers)
[GIN-debug] DELETE /api/v1/users/:id         --> your_project/internal/httptransport.(*handler).DeleteUser-fm (1 handlers)
[GIN-debug] GET    /api/v1/users             --> your_project/internal/httptransport.(*handler).ListUser-fm (1 handlers)
```
6、framework
```
├── config.yaml # 项目配置文件
├── doc # 文档目录
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── viper.go
│   ├── constant.go # 常量定义
│   ├── database
│   │   ├── database.go
│   │   └── mysql.go
│   ├── ent # ent生成目录
│   │   ├── client.go
│   │   ├── ent.go
│   │   ├── enttest
│   │   │   └── enttest.go
│   │   ├── generate.go
│   │   ├── hook
│   │   │   └── hook.go
│   │   ├── migrate
│   │   │   ├── migrate.go
│   │   │   └── schema.go 
│   │   ├── mutation.go
│   │   ├── predicate
│   │   │   └── predicate.go
│   │   ├── runtime
│   │   │   └── runtime.go
│   │   ├── runtime.go
│   │   ├── schema # 数据库schema文件定义
│   │   │   └── yourproject.go
│   │   ├── tx.go
│   │   ├── yourproject
│   │   │   ├── where.go
│   │   │   └── yourproject.go
│   │   ├── yourproject.go
│   │   ├── yourproject_create.go
│   │   ├── yourproject_delete.go
│   │   ├── yourproject_query.go
│   │   └── yourproject_update.go
│   ├── error.go # 错误定义
│   ├── httptransport
│   │   ├── example.go
│   │   ├── gintransport.go # 路由定义
│   │   ├── request # 请求体定义
│   │   └── response # 响应结构定义
│   ├── model # 模型文件
│   │   ├── config.go # 项目配置
│   │   └── pagination.go # 分页结构
│   ├── repository # repo层
│   │   └── repository.go
│   └── service # service层
│       └── service.go
├── main.go
└── script # 项目脚本
├── logrotate.d # 日志切割
│   └── your_project
└── systemd # 使用systemd托管服务
└── your_project.service
```

### More Usage
```shell
wo -h
 
wo is a cli tool for golang backend api with ent.

Usage:
  wo [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        init one golang project
  new         new command for exit init project.

Flags:
  -h, --help   help for wo

Use "wo [command] --help" for more information about a command.

```

### About the Project

`WO` is a cli tool for golang backend api with ent.

### Dependencies
- [entgo](https://github.com/ent/ent)
- [gin-gonic](https://github.com/gin-gonic/gin)