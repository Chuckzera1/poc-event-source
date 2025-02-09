package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func SetupUserRouter(r *gin.RouterGroup) {
	userGroup := r.Group("/user")

	userGroup.POST("", func(context *gin.Context) {
		fmt.Printf("Create User")
	})

	userGroup.GET("/:id", func(context *gin.Context) {
		fmt.Printf("Get User By ID")
	})

	userGroup.PUT("/:id", func(context *gin.Context) {
		fmt.Printf("Update User")
	})
}
