package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	ginbinding "github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"poc-event-source/config"
)

func StartAPI(cfg config.Config, setupRouterFunc func(engine *gin.RouterGroup)) error {
	r := gin.Default()

	if v, ok := ginbinding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("pwd_bytes_max72", func(fl validator.FieldLevel) bool {
			return len([]byte(fl.Field().String())) <= 72
		})
	}

	apiGroup := r.Group("api/v1")
	setupRouterFunc(apiGroup)

	port := fmt.Sprintf(":%v", cfg.APIPort)
	err := r.Run(port)
	if err != nil {
		panic(err)
	}

	return nil
}
