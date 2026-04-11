package routes

import (
	"github.com/ONEST-Network/Job-Manager-Adapter/api/handlers"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	"github.com/gin-gonic/gin"
)

func SchemeRouter(router *gin.RouterGroup, clients *clients.Clients) {
	router.POST("/push", handlers.PushScheme(clients))
	router.GET("/:id/applications", handlers.GetSchemeApplications(clients))
	router.PATCH("/status", handlers.UpdateSchemeStatus(clients))
	router.POST("/:id/apply", handlers.ApplyToScheme(clients))
}
