package scheme

import (
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job"
	jobApp "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job-application"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/scheme"
	schemeapplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/scheme-application"
)

type AddSchemeRequest struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Type          scheme.SchemeType `json:"type"`
	BusinessID    string            `json:"businessId"`
	Eligibility   job.Eligibility   `json:"eligibility"`
	Location      job.Location      `json:"location"`
	FundingAmount float64           `json:"fundingAmount"`
	Category      string            `json:"category"`
}

type UpdateSchemeStatusRequest struct {
	ApplicationID string                         `json:"applicationId"`
	Status        schemeapplication.SchemeApplicationStatus `json:"status"`
}

type JobApplicationRequest struct {
	ApplicantDetails jobApp.ApplicantDetails `json:"applicantDetails"`
}
