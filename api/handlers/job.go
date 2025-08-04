package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ONEST-Network/Job-Manager-Adapter/internal/job"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	jobPayload "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/job"
	"github.com/gin-gonic/gin"
)

// @Summary	Create job
// @Description	Create a job posting
// @Tags Job
// @Accept		json
// @Produce		json
// @Param request body jobPayload.CreateJobRequest true "request body"
// @Success 200
// @Failure 500 {object} string
// @Router	/job/create	[post]
func CreateJob(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload jobPayload.CreateJobRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if err := job.NewJob(clients).CreateJob(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
}

// @Summary	Get job applications
// @Description	Get job applications
// @Tags Job
// @Accept		json
// @Produce		json
// @Param id path string true "Job ID"
// @Success 200 {array} jobPayload.GetJobApplicationsResponse
// @Failure 500 {object} string
// @Router	/job/{id}/applications	[get]
func GetJobApplications(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobID := c.Param("id")

		applications, err := job.NewJob(clients).GetJobApplications(jobID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, applications)
	}
}
