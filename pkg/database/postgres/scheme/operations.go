package scheme

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides access to the scheme storage
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new scheme repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create adds a new scheme to the database
func (r *Repository) Create(ctx context.Context, organizationID int, req CreateSchemeRequest) (*SchemeResponse, error) {
	var scheme Scheme
	query := `
		INSERT INTO schemes (organization_id, scheme_id, title, description, eligibility_criteria, scheme_amount, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'active')
		RETURNING id, organization_id, scheme_id, title, description, eligibility_criteria, scheme_amount, status, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		organizationID,
		req.SchemeID,
		req.Title,
		req.Description,
		req.EligibilityCriteria,
		req.SchemeAmount,
	).Scan(
		&scheme.ID,
		&scheme.OrganizationID,
		&scheme.SchemeID,
		&scheme.Title,
		&scheme.Description,
		&scheme.EligibilityCriteria,
		&scheme.SchemeAmount,
		&scheme.Status,
		&scheme.CreatedAt,
		&scheme.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheme: %v", err)
	}

	return &SchemeResponse{
		ID:                  scheme.ID,
		OrganizationID:      scheme.OrganizationID,
		SchemeID:            scheme.SchemeID,
		Title:               scheme.Title,
		Description:         scheme.Description,
		EligibilityCriteria: scheme.EligibilityCriteria,
		SchemeAmount:        scheme.SchemeAmount,
		Status:              scheme.Status,
		CreatedAt:           scheme.CreatedAt,
		UpdatedAt:           scheme.UpdatedAt,
	}, nil
}

// GetByID retrieves a scheme by its ID
func (r *Repository) GetByID(ctx context.Context, id int) (*SchemeResponse, error) {
	var scheme Scheme
	query := `
		SELECT id, organization_id, scheme_id, title, description, eligibility_criteria, scheme_amount, status, created_at, updated_at
		FROM schemes
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&scheme.ID,
		&scheme.OrganizationID,
		&scheme.SchemeID,
		&scheme.Title,
		&scheme.Description,
		&scheme.EligibilityCriteria,
		&scheme.SchemeAmount,
		&scheme.Status,
		&scheme.CreatedAt,
		&scheme.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get scheme: %v", err)
	}

	return &SchemeResponse{
		ID:                  scheme.ID,
		OrganizationID:      scheme.OrganizationID,
		SchemeID:            scheme.SchemeID,
		Title:               scheme.Title,
		Description:         scheme.Description,
		EligibilityCriteria: scheme.EligibilityCriteria,
		SchemeAmount:        scheme.SchemeAmount,
		Status:              scheme.Status,
		CreatedAt:           scheme.CreatedAt,
		UpdatedAt:           scheme.UpdatedAt,
	}, nil
}

// GetBySchemeID retrieves a scheme by its scheme_id and organization_id
func (r *Repository) GetBySchemeID(ctx context.Context, organizationID int, schemeID string) (*SchemeResponse, error) {
	var scheme Scheme
	query := `
		SELECT id, organization_id, scheme_id, title, description, eligibility_criteria, scheme_amount, status, created_at, updated_at
		FROM schemes
		WHERE organization_id = $1 AND scheme_id = $2
	`

	err := r.db.QueryRow(ctx, query, organizationID, schemeID).Scan(
		&scheme.ID,
		&scheme.OrganizationID,
		&scheme.SchemeID,
		&scheme.Title,
		&scheme.Description,
		&scheme.EligibilityCriteria,
		&scheme.SchemeAmount,
		&scheme.Status,
		&scheme.CreatedAt,
		&scheme.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get scheme: %v", err)
	}

	return &SchemeResponse{
		ID:                  scheme.ID,
		OrganizationID:      scheme.OrganizationID,
		SchemeID:            scheme.SchemeID,
		Title:               scheme.Title,
		Description:         scheme.Description,
		EligibilityCriteria: scheme.EligibilityCriteria,
		SchemeAmount:        scheme.SchemeAmount,
		Status:              scheme.Status,
		CreatedAt:           scheme.CreatedAt,
		UpdatedAt:           scheme.UpdatedAt,
	}, nil
}

// ListByOrganization retrieves all schemes for an organization
func (r *Repository) ListByOrganization(ctx context.Context, organizationID int) ([]*SchemeResponse, error) {
	query := `
		SELECT id, organization_id, scheme_id, title, description, eligibility_criteria, scheme_amount, status, created_at, updated_at
		FROM schemes
		WHERE organization_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list schemes: %v", err)
	}
	defer rows.Close()

	var schemes []*SchemeResponse
	for rows.Next() {
		var scheme Scheme
		if err := rows.Scan(
			&scheme.ID,
			&scheme.OrganizationID,
			&scheme.SchemeID,
			&scheme.Title,
			&scheme.Description,
			&scheme.EligibilityCriteria,
			&scheme.SchemeAmount,
			&scheme.Status,
			&scheme.CreatedAt,
			&scheme.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan scheme: %v", err)
		}

		schemes = append(schemes, &SchemeResponse{
			ID:                  scheme.ID,
			OrganizationID:      scheme.OrganizationID,
			SchemeID:            scheme.SchemeID,
			Title:               scheme.Title,
			Description:         scheme.Description,
			EligibilityCriteria: scheme.EligibilityCriteria,
			SchemeAmount:        scheme.SchemeAmount,
			Status:              scheme.Status,
			CreatedAt:           scheme.CreatedAt,
			UpdatedAt:           scheme.UpdatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return schemes, nil
}

// UpdateStatus updates a scheme's status
func (r *Repository) UpdateStatus(ctx context.Context, id int, status string) (*SchemeResponse, error) {
	var scheme Scheme
	query := `
		UPDATE schemes
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, organization_id, scheme_id, title, description, eligibility_criteria, scheme_amount, status, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, status, time.Now(), id).Scan(
		&scheme.ID,
		&scheme.OrganizationID,
		&scheme.SchemeID,
		&scheme.Title,
		&scheme.Description,
		&scheme.EligibilityCriteria,
		&scheme.SchemeAmount,
		&scheme.Status,
		&scheme.CreatedAt,
		&scheme.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update scheme status: %v", err)
	}

	return &SchemeResponse{
		ID:                  scheme.ID,
		OrganizationID:      scheme.OrganizationID,
		SchemeID:            scheme.SchemeID,
		Title:               scheme.Title,
		Description:         scheme.Description,
		EligibilityCriteria: scheme.EligibilityCriteria,
		SchemeAmount:        scheme.SchemeAmount,
		Status:              scheme.Status,
		CreatedAt:           scheme.CreatedAt,
		UpdatedAt:           scheme.UpdatedAt,
	}, nil
}

// Delete removes a scheme from the database
func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM schemes WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete scheme: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("scheme with ID %d not found", id)
	}

	return nil
}
