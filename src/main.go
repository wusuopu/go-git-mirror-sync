package main

import (
	"app/config"
	"app/initialize"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()
	initialize.Init(e)

	e.Run(":" + config.Config["PORT"].(string))
}
