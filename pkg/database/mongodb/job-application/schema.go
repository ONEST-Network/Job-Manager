package jobapplication

import "time"

type JobApplication struct {
	ID               string               `bson:"id" json:"id"`
	JobID            string               `bson:"job_id" json:"jobId"`
	ApplicantDetails ApplicantDetails     `bson:"applicant_details" json:"applicantDetails"`
	Status           JobApplicationStatus `bson:"status" json:"status"`
	CreatedAt        time.Time            `bson:"created_at" json:"createdAt"`
	UpdatedAt        time.Time            `bson:"updated_at" json:"updatedAt"`
}

type ApplicantDetails struct {
	Name       string     `bson:"name" json:"name"`
	Gender     string     `bson:"gender" json:"gender"`
	Age        int        `bson:"age" json:"age"`
	Experience Experience `bson:"experience" json:"experience"`
	Documents  Documents  `bson:"documents" json:"documents"`
	Phone      string     `bson:"phone" json:"phone"`
	Email      string     `bson:"email" json:"email"`
}

type Documents struct {
	PANCard        *Document `bson:"pan_card" json:"panCard"`
	AadharCard     *Document `bson:"aadhar_card" json:"aadharCard"`
	Passport       *Document `bson:"passport" json:"passport"`
	DrivingLicense *Document `bson:"driving_license" json:"drivingLicense"`
	Resume         *Document `bson:"resume" json:"resume"`
}

type Document struct {
	URL  string `bson:"url" json:"url"`
	Type string `bson:"type" json:"type"`
}

type Experience struct {
	Years int `bson:"years" json:"years"`
}

// JobApplicationStatus represents the status of a job application
type JobApplicationStatus string

const (
	JobApplicationStatusApplicationAccepted  JobApplicationStatus = "APPLICATION_ACCEPTED"
	JobApplicationStatusApplicationRejected  JobApplicationStatus = "APPLICATION_REJECTED"
	JobApplicationStatusAssessmentInProgress JobApplicationStatus = "ASSESSMENT_IN_PROGRESS"
	JobApplicationStatusOfferRejected        JobApplicationStatus = "OFFER_REJECTED"
	JobApplicationStatusOfferAccepted        JobApplicationStatus = "OFFER_ACCEPTED"
	JobApplicationStatusOfferExtended        JobApplicationStatus = "OFFER_EXTENDED"
	JobApplicationStatusCancelled            JobApplicationStatus = "CANCELLED"
)
