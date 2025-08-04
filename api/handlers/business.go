package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ONEST-Network/Job-Manager-Adapter/internal/business"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	businessPayload "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/business"
	"github.com/gin-gonic/gin"
)

// @Summary	Add business
// @Description	Add a business
// @Tags Business
// @Accept		json
// @Produce		json
// @Param request body businessPayload.AddBusinessRequest true "request body"
// @Success 200
// @Failure 500 {object} string
// @Router	/business/add	[post]
func AddBusiness(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload businessPayload.AddBusinessRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		if err := business.NewBusiness(clients).AddBusiness(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
}

// @Summary	List Jobs
// @Description	List jobs for a business
// @Tags Business
// @Accept		json
// @Produce		json
// @Success 200 {array} businessPayload.ListJobsResponse
// @Failure 500 {object} string
// @Router	/business/{id}/jobs	[get]
func ListJobs(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		businessID := c.Param("id")

		jobs, err := business.NewBusiness(clients).ListJobs(businessID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, jobs)
	}
}
