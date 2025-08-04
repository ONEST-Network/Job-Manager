package application

import (
	"encoding/json"
	"time"
)

// Application represents a scholarship application in the database
type Application struct {
	ID                   int             `json:"id"`
	SchemeID             int             `json:"scheme_id"`
	ApplicantID          string          `json:"applicant_id"`
	ApplicantName        string          `json:"applicant_name"`
	ApplicantContact     string          `json:"applicant_contact"`
	ApplicantCredentials json.RawMessage `json:"applicant_credentials"`
	Status               string          `json:"status"`
	IsEligible           bool            `json:"is_eligible"`
	EligibilityDetails   json.RawMessage `json:"eligibility_details"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

// CreateApplicationRequest represents the request to create a new application
type CreateApplicationRequest struct {
	ApplicantID          string          `json:"applicant_id" binding:"required"`
	ApplicantName        string          `json:"applicant_name" binding:"required"`
	ApplicantContact     string          `json:"applicant_contact"`
	ApplicantCredentials json.RawMessage `json:"applicant_credentials" binding:"required"`
}

// UpdateApplicationStatusRequest represents the request to update an application's status
type UpdateApplicationStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// ApplicationResponse represents the response for application operations
type ApplicationResponse struct {
	ID                   int             `json:"id"`
	SchemeID             int             `json:"scheme_id"`
	ApplicantID          string          `json:"applicant_id"`
	ApplicantName        string          `json:"applicant_name"`
	ApplicantContact     string          `json:"applicant_contact"`
	ApplicantCredentials json.RawMessage `json:"applicant_credentials"`
	Status               string          `json:"status"`
	IsEligible           bool            `json:"is_eligible"`
	EligibilityDetails   json.RawMessage `json:"eligibility_details"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}
