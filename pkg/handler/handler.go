package handler

import (
	"mytasks/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/api")
	{
		task := api.Group("/task")
		{
			task.GET("/list", h.getList)
			task.PUT("/update/:id", h.update)
			task.DELETE("/delete/:id", h.delete)
		}
	}

	return router
}
