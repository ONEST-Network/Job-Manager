package routes

import (
	"github.com/ONEST-Network/Job-Manager-Adapter/api/handlers"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	"github.com/gin-gonic/gin"
)

func JobRouter(router *gin.RouterGroup, clients *clients.Clients) {
	router.POST("/create", handlers.CreateJob(clients))
	router.GET("/applications", handlers.GetJobApplications(clients))
}
