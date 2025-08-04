package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides access to the organization storage
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new organization repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create adds a new organization to the database
func (r *Repository) Create(ctx context.Context, req CreateOrganizationRequest) (*OrganizationResponse, error) {
	// Generate a new API key
	apiKey := uuid.New()

	// Insert organization into database
	var org Organization
	query := `
		INSERT INTO organizations (name, api_key, contact_email)
		VALUES ($1, $2, $3)
		RETURNING id, name, api_key, contact_email, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, req.Name, apiKey, req.ContactEmail).Scan(
		&org.ID, &org.Name, &org.APIKey, &org.ContactEmail, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %v", err)
	}

	return &OrganizationResponse{
		ID:           org.ID,
		Name:         org.Name,
		APIKey:       org.APIKey,
		ContactEmail: org.ContactEmail,
		CreatedAt:    org.CreatedAt,
		UpdatedAt:    org.UpdatedAt,
	}, nil
}

// GetByID retrieves an organization by its ID
func (r *Repository) GetByID(ctx context.Context, id int) (*OrganizationResponse, error) {
	var org Organization
	query := `
		SELECT id, name, api_key, contact_email, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&org.ID, &org.Name, &org.APIKey, &org.ContactEmail, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %v", err)
	}

	return &OrganizationResponse{
		ID:           org.ID,
		Name:         org.Name,
		APIKey:       org.APIKey,
		ContactEmail: org.ContactEmail,
		CreatedAt:    org.CreatedAt,
		UpdatedAt:    org.UpdatedAt,
	}, nil
}

// GetByAPIKey retrieves an organization by its API key
func (r *Repository) GetByAPIKey(ctx context.Context, apiKey uuid.UUID) (*OrganizationResponse, error) {
	var org Organization
	query := `
		SELECT id, name, api_key, contact_email, created_at, updated_at
		FROM organizations
		WHERE api_key = $1
	`

	err := r.db.QueryRow(ctx, query, apiKey).Scan(
		&org.ID, &org.Name, &org.APIKey, &org.ContactEmail, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization by API key: %v", err)
	}

	return &OrganizationResponse{
		ID:           org.ID,
		Name:         org.Name,
		APIKey:       org.APIKey,
		ContactEmail: org.ContactEmail,
		CreatedAt:    org.CreatedAt,
		UpdatedAt:    org.UpdatedAt,
	}, nil
}

// List retrieves all organizations
func (r *Repository) List(ctx context.Context) ([]*OrganizationResponse, error) {
	query := `
		SELECT id, name, api_key, contact_email, created_at, updated_at
		FROM organizations
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %v", err)
	}
	defer rows.Close()

	var organizations []*OrganizationResponse
	for rows.Next() {
		var org Organization
		if err := rows.Scan(
			&org.ID, &org.Name, &org.APIKey, &org.ContactEmail, &org.CreatedAt, &org.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan organization: %v", err)
		}

		organizations = append(organizations, &OrganizationResponse{
			ID:           org.ID,
			Name:         org.Name,
			APIKey:       org.APIKey,
			ContactEmail: org.ContactEmail,
			CreatedAt:    org.CreatedAt,
			UpdatedAt:    org.UpdatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return organizations, nil
}

// Update updates an organization's details
func (r *Repository) Update(ctx context.Context, id int, req CreateOrganizationRequest) (*OrganizationResponse, error) {
	var org Organization
	query := `
		UPDATE organizations
		SET name = $1, contact_email = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, name, api_key, contact_email, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, req.Name, req.ContactEmail, time.Now(), id).Scan(
		&org.ID, &org.Name, &org.APIKey, &org.ContactEmail, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update organization: %v", err)
	}

	return &OrganizationResponse{
		ID:           org.ID,
		Name:         org.Name,
		APIKey:       org.APIKey,
		ContactEmail: org.ContactEmail,
		CreatedAt:    org.CreatedAt,
		UpdatedAt:    org.UpdatedAt,
	}, nil
}

// Delete removes an organization from the database
func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM organizations WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("organization with ID %d not found", id)
	}

	return nil
}
