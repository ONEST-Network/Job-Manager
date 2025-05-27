package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ONEST-Network/scheme-manager-adapter/pkg/database/postgres/application"
	"github.com/ONEST-Network/scheme-manager-adapter/pkg/database/postgres/scheme"
	"github.com/ONEST-Network/scheme-manager-adapter/pkg/validator"
	"github.com/gin-gonic/gin"
)

// ApplicationRepositoryInterface defines the interface for application repository
type ApplicationRepositoryInterface interface {
	Create(app *application.Application) error
	GetBySchemeID(schemeID int) ([]application.Application, error)
	GetByID(id int) (*application.Application, error)
	UpdateStatus(id int, status string) error
	Delete(id int) error
}

// SchemeRepositoryInterface defines the interface for scheme repository
type SchemeRepositoryInterface interface {
	GetByID(id int) (*scheme.Scheme, error)
}

// ApplicationHandler handles application-related requests
type ApplicationHandler struct {
	*BaseHandler
	appRepo    *application.Repository
	schemeRepo *scheme.Repository
}

// NewApplicationHandler creates a new application handler
func NewApplicationHandler(base *BaseHandler) *ApplicationHandler {
	return &ApplicationHandler{
		BaseHandler: base,
		appRepo:     application.NewRepository(base.DB),
		schemeRepo:  scheme.NewRepository(base.DB),
	}
}



// Create creates a new application
// @Summary Create a new application
// @Description Create a new application with the given details and validate eligibility
// @Tags applications
// @Accept json
// @Produce json
// @Param scheme_id path int true "Scheme ID"
// @Param application body application.CreateApplicationRequest true "Application details"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /schemes/{scheme_id}/applications [post]
func (h *ApplicationHandler) Create(c *gin.Context) {
	schemeIdStr := c.Param("scheme_id")
	schemeId, err := strconv.Atoi(schemeIdStr)
	if err != nil {
		log.Printf("Error converting scheme_id to int: %v, scheme_id: %s", err, schemeIdStr)
		ErrorResponse(c, http.StatusBadRequest, "Invalid scheme ID", err)
		return
	}

	// Retrieve the scheme details
	scheme, err := h.schemeRepo.GetByID(c.Request.Context(), schemeId)
	if err != nil {
		log.Printf("Error retrieving scheme by ID: %v, scheme_id: %d", err, schemeId)
		ErrorResponse(c, http.StatusNotFound, "Scheme not found", err)
		return
	}

	var req application.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON request: %v", err)
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Validate eligibility
	eligibilityResult, err := validator.ValidateEligibility(
		scheme.EligibilityCriteria,
		req.ApplicantCredentials,
	)
	if err != nil {
		log.Printf("Error validating eligibility: %v, scheme_id: %d", err, schemeId)
		ErrorResponse(c, http.StatusBadRequest, "Failed to validate eligibility", err)
		return
	}

	// Convert eligibility details to JSON
	eligibilityDetails, err := json.Marshal(eligibilityResult.Details)
	if err != nil {
		log.Printf("Error marshaling eligibility details: %v", err)
		ErrorResponse(c, http.StatusInternalServerError, "Failed to process eligibility details", err)
		return
	}

	// Create the application
	app, err := h.appRepo.Create(
		c.Request.Context(),
		schemeId,
		req,
		eligibilityResult.IsEligible,
		eligibilityDetails,
	)
	if err != nil {
		log.Printf("Error creating application: %v, scheme_id: %d", err, schemeId)
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create application", err)
		return
	}

	SuccessResponse(c, http.StatusCreated, "Application created successfully", app)
}

// GetByID retrieves an application by its ID
// @Summary Get an application by ID
// @Description Get application details by ID
// @Tags applications
// @Produce json
// @Param id path int true "Application ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /applications/{id} [get]
func (h *ApplicationHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error converting application id to int: %v, id: %s", err, idStr)
		ErrorResponse(c, http.StatusBadRequest, "Invalid application ID", err)
		return
	}

	app, err := h.appRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Error retrieving application by ID: %v, id: %d", err, id)
		ErrorResponse(c, http.StatusNotFound, "Application not found", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Application retrieved successfully", app)
}

// ListByScheme retrieves all applications for a scheme
// @Summary List applications by scheme
// @Description Get all applications for a specific scheme
// @Tags applications
// @Produce json
// @Param scheme_id path int true "Scheme ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /schemes/{scheme_id}/applications [get]
func (h *ApplicationHandler) ListByScheme(c *gin.Context) {
	schemeIdStr := c.Param("scheme_id")
	schemeId, err := strconv.Atoi(schemeIdStr)
	if err != nil {
		log.Printf("Error converting scheme_id to int: %v, scheme_id: %s", err, schemeIdStr)
		ErrorResponse(c, http.StatusBadRequest, "Invalid scheme ID", err)
		return
	}

	apps, err := h.appRepo.ListByScheme(c.Request.Context(), schemeId)
	if err != nil {
		log.Printf("Error listing applications by scheme: %v, scheme_id: %d", err, schemeId)
		ErrorResponse(c, http.StatusInternalServerError, "Failed to list applications", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Applications retrieved successfully", apps)
}

// UpdateStatus updates an application's status
// @Summary Update application status
// @Description Update the status of an application
// @Tags applications
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Param status body application.UpdateApplicationStatusRequest true "Status update details"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /applications/{id}/status [put]
func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error converting application id to int: %v, id: %s", err, idStr)
		ErrorResponse(c, http.StatusBadRequest, "Invalid application ID", err)
		return
	}

	var req application.UpdateApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON request: %v", err)
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	app, err := h.appRepo.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		log.Printf("Error updating application status: %v, id: %d, status: %s", err, id, req.Status)
		ErrorResponse(c, http.StatusInternalServerError, "Failed to update application status", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Application status updated successfully", app)
}

// Delete deletes an application
// @Summary Delete an application
// @Description Delete an application by ID
// @Tags applications
// @Produce json
// @Param id path int true "Application ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /applications/{id} [delete]
func (h *ApplicationHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error converting application id to int: %v, id: %s", err, idStr)
		ErrorResponse(c, http.StatusBadRequest, "Invalid application ID", err)
		return
	}

	err = h.appRepo.Delete(c.Request.Context(), id)
	if err != nil {
		log.Printf("Error deleting application: %v, id: %d", err, id)
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete application", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Application deleted successfully", nil)
}