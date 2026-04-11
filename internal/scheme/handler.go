package scheme

import (
	"fmt"
	"time"

	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/scheme"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/scheme-application"
	schemePayload "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/scheme"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/utils/random"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type Interface interface {
	AddScheme(payload *schemePayload.AddSchemeRequest) error
	GetSchemeApplications(schemeID string) ([]schemePayload.GetSchemeApplicationsResponse, error)
	UpdateSchemeApplicationStatus(applicationId string, payload *schemePayload.UpdateSchemeStatusRequest) error
	ApplyToScheme(schemeID string, payload *schemePayload.JobApplicationRequest) (string, error)
}

type Scheme struct {
	clients *clients.Clients
}

func NewScheme(clients *clients.Clients) Interface {
	return &Scheme{
		clients: clients,
	}
}

func (s *Scheme) AddScheme(payload *schemePayload.AddSchemeRequest) error {
	logrus.Infof("[Request]: Received request to add a new scheme: %s", payload.Name)

	business, err := s.clients.BusinessClient.GetBusiness(payload.BusinessID)
	if err != nil {
		return fmt.Errorf("failed to get business with id %s, %v", payload.BusinessID, err)
	}

	var schemeObj = &scheme.Scheme{
		ID:            payload.ID,
		Name:          payload.Name,
		Description:   payload.Description,
		Type:          payload.Type,
		Business:      *business,
		Eligibility:   payload.Eligibility,
		Location:      payload.Location,
		FundingAmount: payload.FundingAmount,
		Category:      payload.Category,
	}

	if err := s.clients.SchemeClient.CreateScheme(schemeObj); err != nil {
		logrus.Errorf("Failed to create scheme, %v", err)
		return err
	}

	return nil
}

func (s *Scheme) GetSchemeApplications(schemeID string) ([]schemePayload.GetSchemeApplicationsResponse, error) {
	var response []schemePayload.GetSchemeApplicationsResponse

	logrus.Infof("[Request]: Received request to get applications for scheme %s", schemeID)

	schemeObj, err := s.clients.SchemeClient.GetScheme(schemeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scheme %s, %v", schemeID, err)
	}

	for _, appId := range schemeObj.ApplicationIDs {
		app, err := s.clients.SchemeApplicationClient.GetSchemeApplication(appId)
		if err != nil {
			logrus.Errorf("Failed to get scheme application %s, %v", appId, err)
			continue
		}

		response = append(response, schemePayload.GetSchemeApplicationsResponse{
			ID:               app.ID,
			ApplicantDetails: app.ApplicantDetails,
			Status:           app.Status,
			EligibilityScore: app.EligibilityScore,
		})
	}

	return response, nil
}

func (s *Scheme) UpdateSchemeApplicationStatus(applicationId string, payload *schemePayload.UpdateSchemeStatusRequest) error {
	logrus.Infof("[Request]: Received request to update scheme application %s status", applicationId)

	var (
		query  = bson.D{{Key: "id", Value: applicationId}}
		update = bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: payload.Status}, {Key: "updated_at", Value: time.Now()}}}}
	)

	if err := s.clients.SchemeApplicationClient.UpdateSchemeApplication(query, update); err != nil {
		logrus.Errorf("Failed to update status, %v", err)
		return err
	}

	return nil
}

func (s *Scheme) ApplyToScheme(schemeID string, payload *schemePayload.JobApplicationRequest) (string, error) {
	logrus.Infof("[Request]: Received application for scheme %s from %s", schemeID, payload.ApplicantDetails.Name)

	schemeObj, err := s.clients.SchemeClient.GetScheme(schemeID)
	if err != nil {
		return "", fmt.Errorf("scheme not found, %v", err)
	}

	// Automatic Eligibility Validation
	score, eligible := ValidateEligibility(schemeObj.Eligibility, payload.ApplicantDetails)

	status := schemeapplication.SchemeStatusSubmitted
	if !eligible {
		status = schemeapplication.SchemeStatusIneligible
	} else {
		status = schemeapplication.SchemeStatusEligible
	}

	appID := random.GetRandomString(10)
	app := &schemeapplication.SchemeApplication{
		ID:               appID,
		SchemeID:         schemeID,
		ApplicantDetails: payload.ApplicantDetails,
		Status:           status,
		EligibilityScore: score,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.clients.SchemeApplicationClient.CreateSchemeApplication(app); err != nil {
		return "", err
	}

	// Update Scheme with new Application ID
	query := bson.D{{Key: "id", Value: schemeID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "application_ids", Value: appID}}}}
	if err := s.clients.SchemeClient.UpdateScheme(query, update); err != nil {
		return "", err
	}

	return appID, nil
}
