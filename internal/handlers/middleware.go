package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func normalizeURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			rawURL := c.Query("url")

			if rawURL == "" {
				newErrorResponse(c, http.StatusBadRequest, "missing url parameter")
				return
			}

			normalizedURL, err := url.QueryUnescape(rawURL)
			if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "Invalid URL")
				return
			}

			c.Set("url", normalizedURL)
		}

		c.Next()
	}
}
