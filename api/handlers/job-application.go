package handlers

import (
	"encoding/json"
	"net/http"

	jobApplication "github.com/ONEST-Network/Job-Manager-Adapter/internal/job-application"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	jobapplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/job-application"
	"github.com/gin-gonic/gin"
)

// @Summary	Update job application status
// @Description	Update a job application's status
// @Tags Job Application
// @Accept		json
// @Produce		json
// @Param id path string true "Job Application ID"
// @Param request body jobapplication.UpdateJobApplicationStatusRequest true "request body"
// @Success 200
// @Failure 500 {object} string
// @Router	/job-application/{id}/status	[post]
func UpdateJobApplicationStatus(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		applicationId := c.Param("id")

		var payload jobapplication.UpdateJobApplicationStatusRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if err := jobApplication.NewJobApplication(clients).UpdateJobApplicationStatus(applicationId, &payload); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
}
