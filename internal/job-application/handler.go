package jobapplication

import (
	"fmt"

	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	jobapplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job-application"
	jobApplicationPayload "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/job-application"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type Interface interface {
	UpdateJobApplicationStatus(applicationId string, payload *jobApplicationPayload.UpdateJobApplicationStatusRequest) error
}

type JobApplication struct {
	clients *clients.Clients
}

func NewJobApplication(clients *clients.Clients) Interface {
	return &JobApplication{
		clients: clients,
	}
}

func (j *JobApplication) UpdateJobApplicationStatus(applicationId string, payload *jobApplicationPayload.UpdateJobApplicationStatusRequest) error {
	logrus.Infof("[Request]: Received request to update job application %s status", applicationId)

	jobApplicationStatus, err := getJobApplicationStatusEnum(payload.Status)
	if err != nil {
		logrus.Errorf("Failed while parsing job application status, %v", err)
		return err
	}

	var (
		query  = bson.D{{Key: "id", Value: applicationId}}
		update = bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: jobApplicationStatus}}}}
	)

	if err := j.clients.JobApplicationClient.UpdateJobApplication(query, update); err != nil {
		logrus.Errorf("Failed to update job application %s status, %v", applicationId, err)
		return err
	}

	return nil
}

func getJobApplicationStatusEnum(jobApplicationStatus string) (jobapplication.JobApplicationStatus, error) {
	switch jobApplicationStatus {
	case string(jobapplication.JobApplicationStatusApplicationAccepted):
		return jobapplication.JobApplicationStatusApplicationAccepted, nil
	case string(jobapplication.JobApplicationStatusApplicationRejected):
		return jobapplication.JobApplicationStatusApplicationRejected, nil
	case string(jobapplication.JobApplicationStatusAssessmentInProgress):
		return jobapplication.JobApplicationStatusAssessmentInProgress, nil
	case string(jobapplication.JobApplicationStatusOfferRejected):
		return jobapplication.JobApplicationStatusOfferRejected, nil
	case string(jobapplication.JobApplicationStatusOfferAccepted):
		return jobapplication.JobApplicationStatusOfferAccepted, nil
	case string(jobapplication.JobApplicationStatusOfferExtended):
		return jobapplication.JobApplicationStatusOfferExtended, nil
	case string(jobapplication.JobApplicationStatusCancelled):
		return jobapplication.JobApplicationStatusCancelled, nil
	default:
		return "", fmt.Errorf("invalid job application status %s", jobApplicationStatus)
	}
}
