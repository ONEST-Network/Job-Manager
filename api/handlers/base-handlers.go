package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, APIStatus{Status: "up"})
	}
}

type APIStatus struct {
	Status string `json:"status"`
}
