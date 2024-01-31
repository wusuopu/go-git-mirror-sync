package initialize

import (
	"github.com/joho/godotenv"
)

// 默认加载 .env 文件
func InitEnv(filenames ...string) {
	godotenv.Load(filenames...)
}