package handlers

import (
	"net/http"
	"strconv"

	"github.com/ONEST-Network/scheme-manager-adapter/pkg/database/postgres/scheme"
	"github.com/gin-gonic/gin"
)

// SchemeHandler handles scheme-related requests
type SchemeHandler struct {
	*BaseHandler
	schemeRepo *scheme.Repository
}

// NewSchemeHandler creates a new scheme handler
func NewSchemeHandler(base *BaseHandler) *SchemeHandler {
	return &SchemeHandler{
		BaseHandler: base,
		schemeRepo:  scheme.NewRepository(base.DB),
	}
}

// Create creates a new scheme
// @Summary Create a new scheme
// @Description Create a new scheme with the given details
// @Tags schemes
// @Accept json
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param scheme body scheme.CreateSchemeRequest true "Scheme details"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /organizations/{organization_id}/schemes [post]
func (h *SchemeHandler) Create(c *gin.Context) {
	orgIDStr := c.Param("organization_id")
	orgID, err := strconv.Atoi(orgIDStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid organization ID", err)
		return
	}

	var req scheme.CreateSchemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	scheme, err := h.schemeRepo.Create(c.Request.Context(), orgID, req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create scheme", err)
		return
	}

	SuccessResponse(c, http.StatusCreated, "Scheme created successfully", scheme)
}

// GetBySchemeID retrieves a scheme by its scheme_id and organization_id
// @Summary Get a scheme by scheme_id
// @Description Get scheme details by scheme_id and organization_id
// @Tags schemes
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param scheme_id path string true "Scheme ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /organizations/{organization_id}/schemes/{scheme_id} [get]
func (h *SchemeHandler) GetBySchemeID(c *gin.Context) {
	orgIDStr := c.Param("organization_id")
	orgID, err := strconv.Atoi(orgIDStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid organization ID", err)
		return
	}

	schemeID := c.Param("scheme_id")

	scheme, err := h.schemeRepo.GetBySchemeID(c.Request.Context(), orgID, schemeID)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Scheme not found", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Scheme retrieved successfully", scheme)
}

// ListByOrganization retrieves all schemes for an organization
// @Summary List schemes by organization
// @Description Get all schemes for a specific organization
// @Tags schemes
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /organizations/{organization_id}/schemes [get]
func (h *SchemeHandler) ListByOrganization(c *gin.Context) {
	orgIDStr := c.Param("organization_id")
	orgID, err := strconv.Atoi(orgIDStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid organization ID", err)
		return
	}

	schemes, err := h.schemeRepo.ListByOrganization(c.Request.Context(), orgID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to list schemes", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Schemes retrieved successfully", schemes)
}

// UpdateStatus updates a scheme's status
// @Summary Update scheme status
// @Description Update the status of a scheme
// @Tags schemes
// @Accept json
// @Produce json
// @Param id path int true "Scheme ID"
// @Param status body scheme.UpdateSchemeStatusRequest true "Status update details"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /schemes/{id}/status [put]
func (h *SchemeHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid scheme ID", err)
		return
	}

	var req scheme.UpdateSchemeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	scheme, err := h.schemeRepo.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to update scheme status", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Scheme status updated successfully", scheme)
}

// Delete deletes a scheme
// @Summary Delete a scheme
// @Description Delete a scheme by ID
// @Tags schemes
// @Produce json
// @Param id path int true "Scheme ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /schemes/{id} [delete]
func (h *SchemeHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid scheme ID", err)
		return
	}

	err = h.schemeRepo.Delete(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete scheme", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Scheme deleted successfully", nil)
}
