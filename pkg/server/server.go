package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ONEST-Network/Job-Manager-Adapter/api/middleware"
	"github.com/ONEST-Network/Job-Manager-Adapter/api/routes"
	"github.com/ONEST-Network/Job-Manager-Adapter/docs"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb"
	dbBusiness "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/business"
	dbInitJobApplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/init-job-application"
	dbJob "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job"
	dbJobApplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job-application"
)

func SetupServer(clients *clients.Clients) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(middleware.DefaultStructuredLogger())
	server.Use(gin.Recovery())
	server.Use(middleware.ValidateCors())

	docs.SwaggerInfo.Title = "Job Hub"
	docs.SwaggerInfo.Description = `Job Hub acts as a provider platform (BPP) that hosts the jobs`
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	baseRouter := server.Group("/")
	routes.BaseRouter(baseRouter, clients)
	routes.BecknRouter(baseRouter, clients)

	businessRouter := server.Group("/business")
	routes.BusinessRouter(businessRouter, clients)

	jobRouter := server.Group("/job")
	routes.JobRouter(jobRouter, clients)

	jobApplicationRouter := server.Group("/job-application")
	routes.JobApplicationRouter(jobApplicationRouter, clients)

	return server
}

func InitMongoDB() (*dbBusiness.Dao, *dbJob.Dao, *dbJobApplication.Dao, *dbInitJobApplication.Dao) {
	var err error

	// Initialize mongodb clients
	mongodb.Client, err = mongodb.NewMongoClient()
	if err != nil {
		logrus.Fatalf("[Server]: Failed to connect mongo client, %v", err)
	}

	logrus.Info("[Server]: Connected To MongoDB")

	business := dbBusiness.NewBusinessDao(mongodb.Client.BusinessCollection)
	job := dbJob.NewJobDao(mongodb.Client.JobCollection)
	jobApplication := dbJobApplication.NewJobApplicationDao(mongodb.Client.JobApplicationCollection)
	initJobApplication := dbInitJobApplication.NewInitJobApplicationDao(mongodb.Client.InitJobApplicationCollection)

	return business, job, jobApplication, initJobApplication
}
