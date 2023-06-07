package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) saveURL(c *gin.Context) {
	ctx := context.Background()
	url := c.Query("url")

	if url == "" {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid url").Error())
		return
	}

	code, err := h.service.ShortenURL(ctx, url)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("invalid save url").Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
	})
}

func (h *Handler) redirect(c *gin.Context) {
	ctx := context.Background()
	code := c.Param("code")

	url, err := h.service.GetLongURL(ctx, code)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid save url").Error())
		return
	}

	c.Redirect(http.StatusFound, url)
}
