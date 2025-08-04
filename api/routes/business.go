package routes

import (
	"github.com/ONEST-Network/Job-Manager-Adapter/api/handlers"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	"github.com/gin-gonic/gin"
)

func BusinessRouter(router *gin.RouterGroup, clients *clients.Clients) {
	router.POST("/add", handlers.AddBusiness(clients))
	router.GET("/:id/jobs", handlers.ListJobs(clients))
}
