package scheme

import (
	"encoding/json"
	"time"
)

// Scheme represents a scholarship scheme in the database
type Scheme struct {
	ID                  int             `json:"id"`
	OrganizationID      int             `json:"organization_id"`
	SchemeID            string          `json:"scheme_id"`
	Title               string          `json:"title"`
	Description         string          `json:"description"`
	EligibilityCriteria json.RawMessage `json:"eligibility_criteria"`
	SchemeAmount        float64         `json:"scheme_amount"`
	Status              string          `json:"status"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

// CreateSchemeRequest represents the request to create a new scheme
type CreateSchemeRequest struct {
	SchemeID            string          `json:"scheme_id" binding:"required"`
	Title               string          `json:"title" binding:"required"`
	Description         string          `json:"description"`
	EligibilityCriteria json.RawMessage `json:"eligibility_criteria" binding:"required"`
	SchemeAmount        float64         `json:"scheme_amount"`
}

// UpdateSchemeStatusRequest represents the request to update a scheme's status
type UpdateSchemeStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// SchemeResponse represents the response for scheme operations
type SchemeResponse struct {
	ID                  int             `json:"id"`
	OrganizationID      int             `json:"organization_id"`
	SchemeID            string          `json:"scheme_id"`
	Title               string          `json:"title"`
	Description         string          `json:"description"`
	EligibilityCriteria json.RawMessage `json:"eligibility_criteria"`
	SchemeAmount        float64         `json:"scheme_amount"`
	Status              string          `json:"status"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}
