package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides access to the application storage
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new application repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create adds a new application to the database
func (r *Repository) Create(ctx context.Context, schemeID int, req CreateApplicationRequest, isEligible bool, eligibilityDetails json.RawMessage) (*ApplicationResponse, error) {
	var app Application
	query := `
		INSERT INTO applications (
			scheme_id, applicant_id, applicant_name, applicant_contact, 
			applicant_credentials, status, is_eligible, eligibility_details
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, scheme_id, applicant_id, applicant_name, applicant_contact, 
			applicant_credentials, status, is_eligible, eligibility_details, 
			created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		schemeID,
		req.ApplicantID,
		req.ApplicantName,
		req.ApplicantContact,
		req.ApplicantCredentials,
		"pending",
		isEligible,
		eligibilityDetails,
	).Scan(
		&app.ID,
		&app.SchemeID,
		&app.ApplicantID,
		&app.ApplicantName,
		&app.ApplicantContact,
		&app.ApplicantCredentials,
		&app.Status,
		&app.IsEligible,
		&app.EligibilityDetails,
		&app.CreatedAt,
		&app.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %v", err)
	}

	return &ApplicationResponse{
		ID:                   app.ID,
		SchemeID:             app.SchemeID,
		ApplicantID:          app.ApplicantID,
		ApplicantName:        app.ApplicantName,
		ApplicantContact:     app.ApplicantContact,
		ApplicantCredentials: app.ApplicantCredentials,
		Status:               app.Status,
		IsEligible:           app.IsEligible,
		EligibilityDetails:   app.EligibilityDetails,
		CreatedAt:            app.CreatedAt,
		UpdatedAt:            app.UpdatedAt,
	}, nil
}

// GetByID retrieves an application by its ID
func (r *Repository) GetByID(ctx context.Context, id int) (*ApplicationResponse, error) {
	var app Application
	query := `
		SELECT id, scheme_id, applicant_id, applicant_name, applicant_contact, 
			applicant_credentials, status, is_eligible, eligibility_details, 
			created_at, updated_at
		FROM applications
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&app.ID,
		&app.SchemeID,
		&app.ApplicantID,
		&app.ApplicantName,
		&app.ApplicantContact,
		&app.ApplicantCredentials,
		&app.Status,
		&app.IsEligible,
		&app.EligibilityDetails,
		&app.CreatedAt,
		&app.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %v", err)
	}

	return &ApplicationResponse{
		ID:                   app.ID,
		SchemeID:             app.SchemeID,
		ApplicantID:          app.ApplicantID,
		ApplicantName:        app.ApplicantName,
		ApplicantContact:     app.ApplicantContact,
		ApplicantCredentials: app.ApplicantCredentials,
		Status:               app.Status,
		IsEligible:           app.IsEligible,
		EligibilityDetails:   app.EligibilityDetails,
		CreatedAt:            app.CreatedAt,
		UpdatedAt:            app.UpdatedAt,
	}, nil
}

// ListByScheme retrieves all applications for a scheme
func (r *Repository) ListByScheme(ctx context.Context, schemeID int) ([]*ApplicationResponse, error) {
	query := `
		SELECT id, scheme_id, applicant_id, applicant_name, applicant_contact, 
			applicant_credentials, status, is_eligible, eligibility_details, 
			created_at, updated_at
		FROM applications
		WHERE scheme_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query, schemeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %v", err)
	}
	defer rows.Close()

	var applications []*ApplicationResponse
	for rows.Next() {
		var app Application
		if err := rows.Scan(
			&app.ID,
			&app.SchemeID,
			&app.ApplicantID,
			&app.ApplicantName,
			&app.ApplicantContact,
			&app.ApplicantCredentials,
			&app.Status,
			&app.IsEligible,
			&app.EligibilityDetails,
			&app.CreatedAt,
			&app.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan application: %v", err)
		}

		applications = append(applications, &ApplicationResponse{
			ID:                   app.ID,
			SchemeID:             app.SchemeID,
			ApplicantID:          app.ApplicantID,
			ApplicantName:        app.ApplicantName,
			ApplicantContact:     app.ApplicantContact,
			ApplicantCredentials: app.ApplicantCredentials,
			Status:               app.Status,
			IsEligible:           app.IsEligible,
			EligibilityDetails:   app.EligibilityDetails,
			CreatedAt:            app.CreatedAt,
			UpdatedAt:            app.UpdatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return applications, nil
}

// UpdateStatus updates an application's status
func (r *Repository) UpdateStatus(ctx context.Context, id int, status string) (*ApplicationResponse, error) {
	var app Application
	query := `
		UPDATE applications
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, scheme_id, applicant_id, applicant_name, applicant_contact, 
			applicant_credentials, status, is_eligible, eligibility_details, 
			created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, status, time.Now(), id).Scan(
		&app.ID,
		&app.SchemeID,
		&app.ApplicantID,
		&app.ApplicantName,
		&app.ApplicantContact,
		&app.ApplicantCredentials,
		&app.Status,
		&app.IsEligible,
		&app.EligibilityDetails,
		&app.CreatedAt,
		&app.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update application status: %v", err)
	}

	return &ApplicationResponse{
		ID:                   app.ID,
		SchemeID:             app.SchemeID,
		ApplicantID:          app.ApplicantID,
		ApplicantName:        app.ApplicantName,
		ApplicantContact:     app.ApplicantContact,
		ApplicantCredentials: app.ApplicantCredentials,
		Status:               app.Status,
		IsEligible:           app.IsEligible,
		EligibilityDetails:   app.EligibilityDetails,
		CreatedAt:            app.CreatedAt,
		UpdatedAt:            app.UpdatedAt,
	}, nil
}

// Delete removes an application from the database
func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM applications WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete application: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("application with ID %d not found", id)
	}

	return nil
}
