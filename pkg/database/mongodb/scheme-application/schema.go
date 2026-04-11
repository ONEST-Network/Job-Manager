package schemeapplication

import (
	"time"

	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job-application"
)

type SchemeApplication struct {
	ID               string                 `bson:"id" json:"id"`
	SchemeID         string                 `bson:"scheme_id" json:"schemeId"`
	ApplicantDetails jobapplication.ApplicantDetails `bson:"applicant_details" json:"applicantDetails"`
	Status           SchemeApplicationStatus `bson:"status" json:"status"`
	EligibilityScore float64                `bson:"eligibility_score" json:"eligibilityScore"`
	CreatedAt        time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt        time.Time              `bson:"updated_at" json:"updatedAt"`
}

type SchemeApplicationStatus string

const (
	SchemeStatusSubmitted    SchemeApplicationStatus = "SUBMITTED"
	SchemeStatusEligible     SchemeApplicationStatus = "ELIGIBLE"
	SchemeStatusIneligible   SchemeApplicationStatus = "INELIGIBLE"
	SchemeStatusUnderReview  SchemeApplicationStatus = "UNDER_REVIEW"
	SchemeStatusApproved     SchemeApplicationStatus = "APPROVED"
	SchemeStatusRejected     SchemeApplicationStatus = "REJECTED"
	SchemeStatusDisbursed    SchemeApplicationStatus = "DISBURSED"
)
