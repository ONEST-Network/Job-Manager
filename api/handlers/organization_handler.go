package handlers

import (
	"net/http"
	"strconv"

	"github.com/ONEST-Network/scheme-manager-adapter/pkg/database/organization"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OrganizationHandler handles organization-related requests
type OrganizationHandler struct {
	*BaseHandler
	orgRepo *organization.Repository
}

// NewOrganizationHandler creates a new organization handler
func NewOrganizationHandler(base *BaseHandler) *OrganizationHandler {
	return &OrganizationHandler{
		BaseHandler: base,
		orgRepo:     organization.NewRepository(base.DB),
	}
}

// Create creates a new organization
// @Summary Create a new organization
// @Description Create a new organization with the given details
// @Tags organizations
// @Accept json
// @Produce json
// @Param organization body organization.CreateOrganizationRequest true "Organization details"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /organizations [post]
func (h *OrganizationHandler) Create(c *gin.Context) {
	var req organization.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	org, err := h.orgRepo.Create(c.Request.Context(), req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create organization", err)
		return
	}

	SuccessResponse(c, http.StatusCreated, "Organization created successfully", org)
}

// GetByID retrieves an organization by its ID
// @Summary Get an organization by ID
// @Description Get organization details by ID
// @Tags organizations
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /organizations/{id} [get]
func (h *OrganizationHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid organization ID", err)
		return
	}

	org, err := h.orgRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Organization not found", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Organization retrieved successfully", org)
}

// GetByAPIKey retrieves an organization by its API key
// @Summary Get an organization by API key
// @Description Get organization details by API key
// @Tags organizations
// @Produce json
// @Param api_key path string true "API Key (UUID)"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /organizations/api-key/{api_key} [get]
func (h *OrganizationHandler) GetByAPIKey(c *gin.Context) {
	apiKeyStr := c.Param("api_key")
	apiKey, err := uuid.Parse(apiKeyStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid API key", err)
		return
	}

	org, err := h.orgRepo.GetByAPIKey(c.Request.Context(), apiKey)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Organization not found", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Organization retrieved successfully", org)
}

// List retrieves all organizations
// @Summary List all organizations
// @Description Get a list of all organizations
// @Tags organizations
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /organizations [get]
func (h *OrganizationHandler) List(c *gin.Context) {
	orgs, err := h.orgRepo.List(c.Request.Context())
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to list organizations", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Organizations retrieved successfully", orgs)
}

// Update updates an organization
// @Summary Update an organization
// @Description Update an organization's details
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param organization body organization.CreateOrganizationRequest true "Organization details"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /organizations/{id} [put]
func (h *OrganizationHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid organization ID", err)
		return
	}

	var req organization.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	org, err := h.orgRepo.Update(c.Request.Context(), id, req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to update organization", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Organization updated successfully", org)
}

// Delete deletes an organization
// @Summary Delete an organization
// @Description Delete an organization by ID
// @Tags organizations
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /organizations/{id} [delete]
func (h *OrganizationHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid organization ID", err)
		return
	}

	err = h.orgRepo.Delete(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete organization", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Organization deleted successfully", nil)
}
