# 多个 git 仓库自动同步
## 使用说明


## 开发环境搭建
  * golang 1.21.6

用到的第三方库：
  * Web 框架： github.com/gin-gonic/gin v1.9.1
  * ORM： gorm.io/gorm v1.25.6
  * Migrate: https://github.com/pressly/goose
  * Logger: go.uber.org/zap v1.26.0
  * 单元测试： github.com/stretchr/testify v1.8.4
  * 前端框架： VUE
  * UI库： https://daisyui.com/
  * Jupyter： https://github.com/janpfeifer/gonb

安装 air: `go install github.com/cosmtrek/air@latest`  
安装 goose: `go install github.com/pressly/goose/v3/cmd/goose@v3.18.0`  

配置环境变量： `cp .env.example .env` ；修改对应的配置；针对测试环境： `cp .env.example .evn.test`  
创建数据库表结构： `go run cmd/goose.go [--env development|production|test] up `  
运行开发服务器： `air`  

项目目录结构：
```
src
├── assets/           静态文件
├── cmd/              额外的命令程序
├── config/           解析环境变量之后的配置
├── controllers/      web程序的Controller
├── data/             数据保存的目录
├── di/               di容器
├── initialize/       程序初始化
├── interfaces/       仅定义一些接口，方便单元测试时进行Mock
├── jobs/             定时任务
├── middlewares/      web程序的Middleware
├── migrations/       数据库的 migration
├── models/           数据库的 Model
├── notes/            jupyter notebook
├── routes/           web程序的 Route
├── schemas/          定义一些应用程序的 struct
├── services/         定义一些应用程序的 Service，在执行测试时可以针对性的进行Mock
├── tests/            测试用例
├── tmp/              log保存目录
└── utils/            一些常用功能
```

## 运行测试
执行命令： `go test -v app/tests/...`



