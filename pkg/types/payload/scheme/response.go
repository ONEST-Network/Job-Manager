package scheme

import (
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job-application"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/scheme-application"
)

type GetSchemeApplicationsResponse struct {
	ID               string                               `json:"id"`
	ApplicantDetails jobapplication.ApplicantDetails       `json:"applicantDetails"`
	Status           schemeapplication.SchemeApplicationStatus `json:"status"`
	EligibilityScore float64                              `json:"eligibilityScore"`
}

type ListSchemesResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	FundingAmount float64 `json:"fundingAmount"`
	Category      string  `json:"category"`
}
