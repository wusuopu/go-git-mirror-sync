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

## 运行测试
执行命令： `go test -v app/tests/...`

