package handler

import (
	"mytasks/pkg/cache"
	"mytasks/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	cache    *cache.Cache
}

func NewHandler(services *service.Service, cache *cache.Cache) *Handler {
	return &Handler{services: services, cache: cache}
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
