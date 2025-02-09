package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"poc-event-source/config"
)

func StartAPI(cfg config.Config, setupRouterFunc func(engine *gin.RouterGroup)) error {
	r := gin.Default()

	apiGroup := r.Group("api/v1")
	setupRouterFunc(apiGroup)

	port := fmt.Sprintf(":%v", cfg.APIPort)
	err := r.Run(port)
	if err != nil {
		panic(err)
	}

	return nil
}
