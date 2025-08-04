package business

import (
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/business"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job"
)

type AddBusinessRequest struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Phone          string            `json:"phone"`
	Email          string            `json:"email"`
	PictureURLs    []string          `json:"pictureUrls"`
	Description    string            `json:"description"`
	GSTIndexNumber string            `json:"gstIndexNumber"`
	Location       business.Location `json:"location"`
	Industry       business.Industry `json:"industry"`
}

type ListJobsResponse struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Type           job.JobType     `json:"type"`
	Vacancies      int             `json:"vacancies"`
	SalaryRange    job.SalaryRange `json:"salaryRange"`
	ApplicationIDs []string        `json:"applicationIds"`
	WorkHours      job.WorkHours   `json:"workHours"`
	WorkDays       job.WorkDays    `json:"workDays"`
	Eligibility    job.Eligibility `json:"eligibility"`
	Location       job.Location    `json:"location"`
}
