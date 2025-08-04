package job

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	jobDb "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job"
	jobPayload "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/job"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/utils/random"
)

type Interface interface {
	CreateJob(payload *jobPayload.CreateJobRequest) error
	GetJobApplications(jobID string) ([]jobPayload.GetJobApplicationsResponse, error)
}

type Job struct {
	clients *clients.Clients
}

func NewJob(clients *clients.Clients) Interface {
	return &Job{
		clients: clients,
	}
}

func (j *Job) CreateJob(payload *jobPayload.CreateJobRequest) error {
	logrus.Infof("[Request]: Received request to create a new job for business: %s", payload.BusinessID)

	business, err := j.clients.BusinessClient.GetBusiness(payload.BusinessID)
	if err != nil {
		return fmt.Errorf("failed to get business with id %s, %v", payload.BusinessID, err)
	}

	var job = &jobDb.Job{
		ID:          random.GetRandomString(7),
		Name:        payload.Name,
		Description: payload.Description,
		Type:        payload.Type,
		Vacancies:   payload.Vacancies,
		SalaryRange: payload.SalaryRange,
		WorkHours:   payload.WorkHours,
		WorkDays:    payload.WorkDays,
		Eligibility: payload.Eligibility,
		Location:    payload.Location,
		Business:    *business,
	}

	if err := j.clients.JobClient.CreateJob(job); err != nil {
		logrus.Errorf("Failed to create job, %v", err)
		return err
	}

	return nil
}

func (j *Job) GetJobApplications(jobID string) ([]jobPayload.GetJobApplicationsResponse, error) {
	var jobApplicationsResponse []jobPayload.GetJobApplicationsResponse

	logrus.Infof("[Request]: Received request to get applications for %s job", jobID)

	job, err := j.clients.JobClient.GetJob(jobID)
	if err != nil {
		logrus.Errorf("Failed to get job %s, %v", jobID, err)
		return nil, fmt.Errorf("failed to get job %s, %v", jobID, err)
	}

	for _, applicationId := range job.ApplicationIDs {
		jobApplication, err := j.clients.JobApplicationClient.GetJobApplication(applicationId)
		if err != nil {
			logrus.Errorf("Failed to get job application %s, %v", applicationId, err)
			return nil, fmt.Errorf("failed to get job application %s, %v", applicationId, err)
		}

		jobApplicationsResponse = append(jobApplicationsResponse, jobPayload.GetJobApplicationsResponse{
			ID:               applicationId,
			ApplicantDetails: jobApplication.ApplicantDetails,
			Status:           jobApplication.Status,
		})
	}

	return jobApplicationsResponse, nil
}
