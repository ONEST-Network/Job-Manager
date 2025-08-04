package postgres

import (
	"context"
	"fmt"

	"github.com/ONEST-Network/scheme-manager-adapter/pkg/config"
	// "github.com/ONEST-Network/scheme-manager-adapter/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDB initializes the database connection
func InitDB(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	// Use the PostgresConnectionString method from DatabaseConfig
	connString := cfg.Database.PostgresConnectionString()

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	// Set connection pool settings from config
	poolConfig.MaxConns = int32(cfg.Database.MaxConns)

	// Connect to database
	dbpool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	// Ping database to verify connection
	if err := dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	// Initialize database tables
	if err := initTables(ctx, dbpool); err != nil {
		return nil, fmt.Errorf("failed to initialize database tables: %v", err)
	}

	return dbpool, nil
}

// initTables creates the necessary tables if they don't exist
func initTables(ctx context.Context, db *pgxpool.Pool) error {
	// SQL for creating organizations table
	createOrganizationsTable := `
	CREATE TABLE IF NOT EXISTS organizations (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		api_key UUID NOT NULL UNIQUE,
		contact_email VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	// SQL for creating schemes table
	createSchemesTable := `
	CREATE TABLE IF NOT EXISTS schemes (
		id SERIAL PRIMARY KEY,
		organization_id INTEGER REFERENCES organizations(id),
		scheme_id VARCHAR(255) NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		eligibility_criteria JSONB,
		scheme_amount NUMERIC(10,2),
		status VARCHAR(50) DEFAULT 'active',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(organization_id, scheme_id)
	)`

	// SQL for creating applications table
	createApplicationsTable := `
	CREATE TABLE IF NOT EXISTS applications (
		id SERIAL PRIMARY KEY,
		scheme_id INTEGER REFERENCES schemes(id),
		applicant_id VARCHAR(255) NOT NULL,
		applicant_name VARCHAR(255),
		applicant_contact VARCHAR(255),
		applicant_credentials JSONB,
		status VARCHAR(50) DEFAULT 'pending',
		is_eligible BOOLEAN DEFAULT FALSE,
		eligibility_details JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(scheme_id, applicant_id)
	)`

	// Execute the SQL statements
	if _, err := db.Exec(ctx, createOrganizationsTable); err != nil {
		return fmt.Errorf("failed to create organizations table: %v", err)
	}

	if _, err := db.Exec(ctx, createSchemesTable); err != nil {
		return fmt.Errorf("failed to create schemes table: %v", err)
	}

	if _, err := db.Exec(ctx, createApplicationsTable); err != nil {
		return fmt.Errorf("failed to create applications table: %v", err)
	}

	return nil
}
