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
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing URL parameter"})
				return
			}

			normalizedURL, err := url.QueryUnescape(rawURL)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
				return
			}

			c.Set("url", normalizedURL)
		}

		c.Next()
	}
}
