package routes

import (
	"github.com/ONEST-Network/Job-Manager-Adapter/api/handlers"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	"github.com/gin-gonic/gin"
)

func JobApplicationRouter(router *gin.RouterGroup, clients *clients.Clients) {
	router.POST("/:id/status", handlers.UpdateJobApplicationStatus(clients))
}
