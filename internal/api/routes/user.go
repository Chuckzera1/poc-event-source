package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"poc-event-source/internal/application/dto"
	usecaseuser "poc-event-source/internal/application/usecase/user"
)

type UserHandler struct {
	createUser usecaseuser.CreateUserUseCase
}

func NewUserHandler(cu usecaseuser.CreateUserUseCase) *UserHandler {
	return &UserHandler{createUser: cu}
}

func (h *UserHandler) SetupUserRouter(r *gin.RouterGroup) {
	g := r.Group("/user")
	g.POST("", h.create)
	g.GET("/:id", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, nil) })
	g.PUT("/:id", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, nil) })
}

func (h *UserHandler) create(c *gin.Context) {
	var req dto.CreateUserReqDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.createUser.Execute(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "event published"})
}
