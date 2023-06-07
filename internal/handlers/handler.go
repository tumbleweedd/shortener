package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tumbleweedd/shortener/internal/services"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/a", h.saveURL)
	router.GET("/s/:code", h.redirect)

	return router
}
