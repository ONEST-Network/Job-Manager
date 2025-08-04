package job

import (
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job"
	jobapplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job-application"
)

type CreateJobRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Type        job.JobType     `json:"type"`
	Vacancies   int             `json:"vacancies"`
	SalaryRange job.SalaryRange `json:"salaryRange"`
	WorkHours   job.WorkHours   `json:"workHours"`
	WorkDays    job.WorkDays    `json:"workDays"`
	Eligibility job.Eligibility `json:"eligibility"`
	Location    job.Location    `json:"location"`
	BusinessID  string          `json:"businessId"`
}

type GetJobApplicationsResponse struct {
	ID               string                              `json:"id"`
	ApplicantDetails jobapplication.ApplicantDetails     `json:"applicantDetails"`
	Status           jobapplication.JobApplicationStatus `json:"status"`
}
