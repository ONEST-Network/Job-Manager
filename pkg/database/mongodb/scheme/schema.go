package scheme

import (
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/business"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job"
)

// Scheme represents a scheme in the database
type Scheme struct {
	ID             string            `bson:"id" json:"id"`
	Name           string            `bson:"name" json:"name"`
	Description    string            `bson:"description" json:"description"`
	Type           SchemeType        `bson:"type" json:"type"`
	ApplicationIDs []string          `bson:"application_ids" json:"applicationIds"`
	Business       business.Business `bson:"business" json:"business"`
	Eligibility    job.Eligibility   `bson:"eligibility" json:"eligibility"`
	Location       job.Location      `bson:"location" json:"location"`
	FundingAmount  float64           `bson:"funding_amount" json:"fundingAmount"`
	Category       string            `bson:"category" json:"category"`
}

type SchemeType string

const (
	SchemeTypeScholarship SchemeType = "scholarship"
	SchemeTypeGrant       SchemeType = "grant"
	SchemeTypeTraining    SchemeType = "training"
	SchemeTypeIncubation  SchemeType = "incubation"
)
