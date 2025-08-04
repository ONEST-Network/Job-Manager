package organization

import (
	"time"

	"github.com/google/uuid"
)

// Organization represents an organization record in the database
type Organization struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	APIKey       uuid.UUID `json:"api_key"`
	ContactEmail string    `json:"contact_email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateOrganizationRequest represents the request to create a new organization
type CreateOrganizationRequest struct {
	Name         string `json:"name" binding:"required"`
	ContactEmail string `json:"contact_email" binding:"required,email"`
}

// OrganizationResponse represents the response for organization operations
type OrganizationResponse struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	APIKey       uuid.UUID `json:"api_key"`
	ContactEmail string    `json:"contact_email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
