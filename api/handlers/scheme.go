package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ONEST-Network/Job-Manager-Adapter/internal/scheme"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	schemePayload "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/scheme"
	"github.com/gin-gonic/gin"
)

// @Summary	Push Scheme
// @Description	Add a new scheme to the manager
// @Tags Scheme
// @Accept		json
// @Produce		json
// @Param request body schemePayload.AddSchemeRequest true "request body"
// @Success 200
// @Failure 500 {object} string
// @Router	/scheme/push	[post]
func PushScheme(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload schemePayload.AddSchemeRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if err := scheme.NewScheme(clients).AddScheme(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Scheme added successfully"})
	}
}

// @Summary	Get scheme applications
// @Description	Get applications for a specific scheme
// @Tags Scheme
// @Accept		json
// @Produce		json
// @Param id path string true "Scheme ID"
// @Success 200 {array} schemePayload.GetSchemeApplicationsResponse
// @Failure 500 {object} string
// @Router	/scheme/{id}/applications	[get]
func GetSchemeApplications(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		schemeID := c.Param("id")

		applications, err := scheme.NewScheme(clients).GetSchemeApplications(schemeID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, applications)
	}
}

// @Summary	Update scheme application status
// @Description	Update the status of a scheme application
// @Tags Scheme
// @Accept		json
// @Produce		json
// @Param request body schemePayload.UpdateSchemeStatusRequest true "request body"
// @Success 200
// @Failure 500 {object} string
// @Router	/scheme/status	[patch]
func UpdateSchemeStatus(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload schemePayload.UpdateSchemeStatusRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if err := scheme.NewScheme(clients).UpdateSchemeApplicationStatus(payload.ApplicationID, &payload); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
	}
}

// @Summary	Apply to Scheme
// @Description	Submit an application for a scheme (Triggers auto-eligibility)
// @Tags Scheme
// @Accept		json
// @Produce		json
// @Param id path string true "Scheme ID"
// @Param request body schemePayload.JobApplicationRequest true "request body"
// @Success 200 {object} string "Application ID"
// @Failure 500 {object} string
// @Router	/scheme/{id}/apply [post]
func ApplyToScheme(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		schemeID := c.Param("id")
		var payload schemePayload.JobApplicationRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		appID, err := scheme.NewScheme(clients).ApplyToScheme(schemeID, &payload)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"applicationId": appID})
	}
}
