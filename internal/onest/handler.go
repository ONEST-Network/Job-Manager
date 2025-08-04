package onest

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/builders/onest"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	dbInitJobApplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/init-job-application"
	dbJobApplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job-application"

	searchrequest "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/search/request"
	searchrequestack "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/search/request-ack"
	searchresponseack "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/search/response-ack"

	selectrequest "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/select/request"
	selectrequestack "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/select/request-ack"
	selectresponseack "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/select/response-ack"

	initrequest "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/init/request"
	initrequestack "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/init/request-ack"
	initresponseack "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/init/response-ack"

	confirmrequest "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/confirm/request"
	confirmrequestack "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/confirm/request-ack"
	confirmresponseack "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/confirm/response-ack"

	statusrequest "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/status/request"
	statusrequestack "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/status/request-ack"
	statusresponseack "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/status/response-ack"

	cancelrequest "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/cancel/request"
	cancelrequestack "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/cancel/request-ack"
	cancelresponseack "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/cancel/response-ack"
)

type Interface interface {
	// search api handlers
	SendJobsAck(body io.ReadCloser) (*searchrequest.SearchRequest, *searchrequestack.SearchRequestAck)
	SendJobs(payload *searchrequest.SearchRequest)
	// select api handlers
	SendJobFulfillmentAck(body io.ReadCloser) (*selectrequest.SelectRequest, *selectrequestack.SelectRequestAck)
	SendJobFulfillment(payload *selectrequest.SelectRequest)
	// init api handlers
	InitializeJobApplicationAck(body io.ReadCloser) (*initrequest.InitRequest, *initrequestack.InitRequestAck)
	InitializeJobApplication(payload *initrequest.InitRequest)
	// confirm api handlers
	ConfirmJobApplicationAck(body io.ReadCloser) (*confirmrequest.ConfirmRequest, *dbInitJobApplication.InitJobApplication, *confirmrequestack.ConfirmRequestAck)
	ConfirmJobApplication(payload *confirmrequest.ConfirmRequest, initJobApplication *dbInitJobApplication.InitJobApplication)
	// status api handlers
	JobApplicationStatusAck(body io.ReadCloser) (*statusrequest.StatusRequest, *statusrequestack.StatusRequestAck)
	JobApplicationStatus(payload *statusrequest.StatusRequest)
	// cancel api handlers
	WithdrawJobApplicationAck(body io.ReadCloser) (*cancelrequest.CancelRequest, *cancelrequestack.CancelRequestAck)
	WithdrawJobApplication(payload *cancelrequest.CancelRequest)
}

type Onest struct {
	clients *clients.Clients
}

func NewOnestClient(clients *clients.Clients) Interface {
	return &Onest{
		clients: clients,
	}
}

func (j *Onest) SendJobsAck(body io.ReadCloser) (*searchrequest.SearchRequest, *searchrequestack.SearchRequestAck) {
	var (
		payload      searchrequest.SearchRequest
		payloadError *searchrequestack.Error
		status       = "ACK"
	)

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		payloadError = &searchrequestack.Error{
			Code:    "10000",
			Paths:   "",
			Message: err.Error(),
		}
	}

	if payloadError != nil {
		status = "NACK"
	}

	return &payload, &searchrequestack.SearchRequestAck{
		Message: searchrequestack.Message{
			Ack: searchrequestack.Ack{
				Status: status,
			},
		},
		Error: payloadError,
	}
}

func (j *Onest) SendJobs(payload *searchrequest.SearchRequest) {
	jobs, err := j.clients.JobClient.ListJobs(getSearchFilter(payload))
	if err != nil {
		logrus.Errorf("Failed to list jobs, %v", err)
		return
	}

	response, err := onest.BuildSearchJobsResponse(j.clients, payload, jobs)
	if err != nil {
		logrus.Errorf("Failed to build list jobs, %v", err)
		return
	}

	var searchResponseAck searchresponseack.SearchResponseAck
	if err := j.clients.ApiClient.ApiCall(response, payload.Context.BapURI+"/on_search", &searchResponseAck, "POST"); err != nil {
		logrus.Errorf("Failed to send jobs search response, %v", err)
		return
	}

	if searchResponseAck.Error.Message != "" {
		logrus.Errorf("Received error while sending jobs, %v", searchResponseAck.Error.Message)
		return
	}
}

func (j *Onest) SendJobFulfillmentAck(body io.ReadCloser) (*selectrequest.SelectRequest, *selectrequestack.SelectRequestAck) {
	var (
		payload  selectrequest.SelectRequest
		getError = func(message, paths, code string) *selectrequestack.SelectRequestAck {
			if code == "" {
				code = "10000"
			}

			return &selectrequestack.SelectRequestAck{
				Message: selectrequestack.Message{
					Ack: selectrequestack.Ack{
						Status: "NACK",
					},
				},
				Error: &selectrequestack.Error{
					Code:    code,
					Paths:   paths,
					Message: message,
				},
			}
		}
	)

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return nil, getError(err.Error(), "", "")
	}

	if payload.Message.Order.Items == nil {
		return nil, getError("No items found", ".message.order.items", "30004")
	}

	jobs, err := j.clients.JobClient.ListJobs(bson.D{{Key: "id", Value: payload.Message.Order.Items[0].ID}})
	if err != nil {
		return nil, getError(err.Error(), "", "")
	}

	if jobs == nil {
		return nil, getError("No job found for id: "+payload.Message.Order.Items[0].ID, "", "30004")
	}

	if jobs[0].Vacancies == 0 {
		return nil, getError("No vacancies available for job: "+jobs[0].ID, "", "40002")
	}

	return &payload, &selectrequestack.SelectRequestAck{
		Message: selectrequestack.Message{
			Ack: selectrequestack.Ack{
				Status: "ACK",
			},
		},
	}
}

func (j *Onest) SendJobFulfillment(payload *selectrequest.SelectRequest) {
	response := onest.BuildSendJobFulfillmentResponse(payload)

	var selectResponseAck selectresponseack.SelectResponseAck
	if err := j.clients.ApiClient.ApiCall(response, payload.Context.BapURI+"/on_select", &selectResponseAck, "POST"); err != nil {
		logrus.Errorf("Failed to send job fulfillment response, %v", err)
		return
	}

	if selectResponseAck.Error.Message != "" {
		logrus.Errorf("Received error while sending job fulfillment, %v", selectResponseAck.Error.Message)
		return
	}
}

func (j *Onest) InitializeJobApplicationAck(body io.ReadCloser) (*initrequest.InitRequest, *initrequestack.InitRequestAck) {
	var (
		payload  initrequest.InitRequest
		getError = func(message, paths, code string) *initrequestack.InitRequestAck {
			if code == "" {
				code = "10000"
			}

			return &initrequestack.InitRequestAck{
				Message: initrequestack.Message{
					Ack: initrequestack.Ack{
						Status: "NACK",
					},
				},
				Error: &initrequestack.Error{
					Code:    code,
					Paths:   paths,
					Message: message,
				},
			}
		}
	)

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return nil, getError(err.Error(), "", "")
	}

	if payload.Message.Order.Fulfillments == nil {
		return nil, getError("no fulfillments found", ".message.order.fulfillments", "")
	}

	if payload.Message.Order.Items == nil {
		return nil, getError("no items found", ".message.order.items", "")
	}

	return &payload, &initrequestack.InitRequestAck{
		Message: initrequestack.Message{
			Ack: initrequestack.Ack{
				Status: "ACK",
			},
		},
		Error: nil,
	}
}

func (j *Onest) InitializeJobApplication(payload *initrequest.InitRequest) {
	age, err := strconv.Atoi(payload.Message.Order.Fulfillments[0].Customer.Person.Age)
	if err != nil {
		logrus.Errorf("Failed to convert age to int, %v", err)
	}

	experience, err := getExeperience(payload)
	if err != nil {
		logrus.Errorf("Failed to get experience, %v", err)
	}

	initJobApplication := dbInitJobApplication.InitJobApplication{
		TransactionID: payload.Context.TransactionID,
		JobID:         payload.Message.Order.Items[0].ID,
		CreatedAt:     time.Now(),
		ApplicantDetails: dbInitJobApplication.ApplicantDetails{
			Name:   payload.Message.Order.Fulfillments[0].Customer.Person.Name,
			Gender: payload.Message.Order.Fulfillments[0].Customer.Person.Gender,
			Age:    age,
			Experience: dbInitJobApplication.Experience{
				Years: experience,
			},
			Documents: getCredentials(payload),
			Phone:     payload.Message.Order.Fulfillments[0].Customer.Contact.Phone,
			Email:     payload.Message.Order.Fulfillments[0].Customer.Contact.Email,
		},
	}

	if err := j.clients.InitJobApplicationClient.CreateInitJobApplication(&initJobApplication); err != nil {
		logrus.Errorf("Failed to create init job application, %v", err)
		return
	}

	response := onest.BuildInitializeJobApplicationResponse(payload)

	var initResponseAck initresponseack.InitResponseAck
	if err := j.clients.ApiClient.ApiCall(response, payload.Context.BapURI+"/on_init", &initResponseAck, "POST"); err != nil {
		logrus.Errorf("Failed to send init job application response, %v", err)
		return
	}

	if initResponseAck.Error.Message != "" {
		logrus.Errorf("Received error while sending job application init response, %v", initResponseAck.Error.Message)
		return
	}
}

func (j *Onest) ConfirmJobApplicationAck(body io.ReadCloser) (*confirmrequest.ConfirmRequest, *dbInitJobApplication.InitJobApplication, *confirmrequestack.ConfirmRequestAck) {
	var (
		payload  confirmrequest.ConfirmRequest
		getError = func(message, paths, code string) *confirmrequestack.ConfirmRequestAck {
			if code == "" {
				code = "10000"
			}

			return &confirmrequestack.ConfirmRequestAck{
				Message: confirmrequestack.Message{
					Ack: confirmrequestack.Ack{
						Status: "NACK",
					},
				},
				Error: &confirmrequestack.Error{
					Code:    code,
					Paths:   paths,
					Message: message,
				},
			}
		}
	)

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return nil, nil, getError(err.Error(), "", "")
	}

	if payload.Message.Order.Fulfillments == nil {
		logrus.Errorf("No fulfillments found")
		return nil, nil, getError("no fulfillments found", ".message.order.fulfillments", "30004")
	}

	if payload.Message.Order.Items == nil {
		logrus.Errorf("No items found")
		return nil, nil, getError("no items found", ".message.order.items", "30004")
	}

	initJobApplication, err := j.clients.InitJobApplicationClient.GetInitJobApplication(payload.Context.TransactionID)
	if err != nil {
		logrus.Errorf("No init job application found for %s transaction-id, %v", payload.Context.TransactionID, err)
		return nil, nil, getError("no init job application found for the given transaction-id", "", "30004")
	}

	return &payload, initJobApplication, &confirmrequestack.ConfirmRequestAck{
		Message: confirmrequestack.Message{
			Ack: confirmrequestack.Ack{
				Status: "ACK",
			},
		},
		Error: nil,
	}
}

func (j *Onest) ConfirmJobApplication(payload *confirmrequest.ConfirmRequest, initJobApplication *dbInitJobApplication.InitJobApplication) {
	if err := j.clients.InitJobApplicationClient.DeleteInitJobApplication(payload.Context.TransactionID); err != nil {
		logrus.Errorf("Failed to delete init job application for %s transaction-id, %v", payload.Context.TransactionID, err)
	}

	if err := j.clients.JobApplicationClient.CreateJobApplication(&dbJobApplication.JobApplication{
		ID:    payload.Message.Order.ID,
		JobID: payload.Message.Order.Items[0].ID,
		ApplicantDetails: dbJobApplication.ApplicantDetails{
			Name:   initJobApplication.ApplicantDetails.Name,
			Gender: initJobApplication.ApplicantDetails.Gender,
			Age:    initJobApplication.ApplicantDetails.Age,
			Experience: dbJobApplication.Experience{
				Years: initJobApplication.ApplicantDetails.Experience.Years,
			},
			Documents: getJobApplicationDocuments(initJobApplication),
			Phone:     initJobApplication.ApplicantDetails.Phone,
			Email:     initJobApplication.ApplicantDetails.Email,
		},
		Status:    dbJobApplication.JobApplicationStatusApplicationAccepted,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		logrus.Errorf("Failed to create %s job application, %v", payload.Message.Order.ID, err)
		return
	}

	var (
		query  = bson.D{{Key: "id", Value: payload.Message.Order.Items[0].ID}}
		update = bson.D{{Key: "$inc", Value: bson.D{{Key: "vacancies", Value: -1}}}}
	)

	if err := j.clients.JobClient.UpdateJob(query, update); err != nil {
		logrus.Errorf("Failed to update %s job vacancies, %v", payload.Message.Order.Items[0].ID, err)
		return
	}

	response := onest.BuildConfirmJobApplicationResponse(payload)

	var initResponseAck confirmresponseack.ConfirmResponseAck
	if err := j.clients.ApiClient.ApiCall(response, payload.Context.BapURI+"/on_confirm", &initResponseAck, "POST"); err != nil {
		logrus.Errorf("Failed to send jobs, %v", err)
		return
	}

	if initResponseAck.Error.Message != "" {
		logrus.Errorf("Received error while sending job application creation response, %v", initResponseAck.Error.Message)
		return
	}
}

func (j *Onest) JobApplicationStatusAck(body io.ReadCloser) (*statusrequest.StatusRequest, *statusrequestack.StatusRequestAck) {
	var (
		payload  statusrequest.StatusRequest
		getError = func(message, paths, code string) *statusrequestack.StatusRequestAck {
			if code == "" {
				code = "10000"
			}

			return &statusrequestack.StatusRequestAck{
				Message: statusrequestack.Message{
					Ack: statusrequestack.Ack{
						Status: "NACK",
					},
				},
				Error: &statusrequestack.Error{
					Code:    code,
					Paths:   paths,
					Message: message,
				},
			}
		}
	)

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return nil, getError(err.Error(), "", "")
	}

	return &payload, &statusrequestack.StatusRequestAck{
		Message: statusrequestack.Message{
			Ack: statusrequestack.Ack{
				Status: "ACK",
			},
		},
		Error: nil,
	}
}

func (j *Onest) JobApplicationStatus(payload *statusrequest.StatusRequest) {
	jobApplication, err := j.clients.JobApplicationClient.GetJobApplication(payload.Message.Order.ID)
	if err != nil {
		logrus.Errorf("Failed to get %s job application, %v", payload.Message.Order.ID, err)
		return
	}

	response := onest.BuildJobApplicationStatusResponse(payload, jobApplication)

	var statusResponseAck statusresponseack.StatusResponseAck
	if err := j.clients.ApiClient.ApiCall(response, payload.Context.BapURI+"/on_status", &statusResponseAck, "POST"); err != nil {
		logrus.Errorf("Failed to send job application status response, %v", err)
		return
	}

	if statusResponseAck.Error.Message != "" {
		logrus.Errorf("Received error while sending job application status response, %v", statusResponseAck.Error.Message)
		return
	}
}

func (j *Onest) WithdrawJobApplicationAck(body io.ReadCloser) (*cancelrequest.CancelRequest, *cancelrequestack.CancelRequestAck) {
	var (
		payload  cancelrequest.CancelRequest
		getError = func(message, paths, code string) *cancelrequestack.CancelRequestAck {
			if code == "" {
				code = "10000"
			}

			return &cancelrequestack.CancelRequestAck{
				Message: cancelrequestack.Message{
					Ack: cancelrequestack.Ack{
						Status: "NACK",
					},
				},
				Error: &cancelrequestack.Error{
					Code:    code,
					Paths:   paths,
					Message: message,
				},
			}
		}
	)

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return nil, getError(err.Error(), "", "")
	}

	return &payload, &cancelrequestack.CancelRequestAck{
		Message: cancelrequestack.Message{
			Ack: cancelrequestack.Ack{
				Status: "ACK",
			},
		},
		Error: nil,
	}
}

func (j *Onest) WithdrawJobApplication(payload *cancelrequest.CancelRequest) {
	var (
		query  = bson.D{{Key: "id", Value: payload.Message.OrderID}}
		update = bson.D{{Key: "$set", Value: bson.D{
			{Key: "status", Value: dbJobApplication.JobApplicationStatusCancelled},
		}}}
	)

	jobApplication, err := j.clients.JobApplicationClient.UpdateJobApplicationAndReturnDocument(query, update)
	if err != nil {
		logrus.Errorf("Failed to update %s job application as withdrawn, %v", payload.Message.OrderID, err)
		return
	}

	query = bson.D{{Key: "id", Value: jobApplication.JobID}}
	update = bson.D{{Key: "$inc", Value: bson.D{{Key: "vacancies", Value: 1}}}}

	if err := j.clients.JobClient.UpdateJob(query, update); err != nil {
		logrus.Errorf("Failed to update %s job vacancies, %v", jobApplication.JobID, err)
		return
	}

	response := onest.BuildWithdrawJobApplicationResponse(payload, jobApplication)

	var cancelResponseAck cancelresponseack.CancelResponseAck
	if err := j.clients.ApiClient.ApiCall(response, payload.Context.BapURI+"/on_cancel", &cancelResponseAck, "POST"); err != nil {
		logrus.Errorf("Failed to send job application withdrawal response, %v", err)
		return
	}

	if cancelResponseAck.Error.Message != "" {
		logrus.Errorf("Received error while sending job application withdrawal response, %v", cancelResponseAck.Error.Message)
		return
	}
}

func getExeperience(payload *initrequest.InitRequest) (int, error) {
	for _, tag := range payload.Message.Order.Fulfillments[0].Customer.Person.Tags {
		if tag.Descriptor.Code == "WORK_EXPERIENCE" {
			if len(tag.List) == 0 {
				return 0, nil
			}
			return extractYears(tag.List[0].Value)
		}
	}

	return 0, nil
}

func extractYears(duration string) (int, error) {
	re := regexp.MustCompile(`P(\d+)Y`)
	matches := re.FindStringSubmatch(duration)
	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid duration format")
	}
	years, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}
	return years, nil
}

func getCredentials(payload *initrequest.InitRequest) dbInitJobApplication.Documents {
	var documents dbInitJobApplication.Documents

	for _, cred := range payload.Message.Order.Fulfillments[0].Customer.Person.Creds {
		switch cred.Descriptor.Name {
		case "PAN_CARD":
			documents.PANCard = &dbInitJobApplication.Document{
				URL:  cred.URL,
				Type: cred.Type,
			}
		case "AADHAAR_CARD":
			documents.AadharCard = &dbInitJobApplication.Document{
				URL:  cred.URL,
				Type: cred.Type,
			}
		case "PASSPORT":
			documents.Passport = &dbInitJobApplication.Document{
				URL:  cred.URL,
				Type: cred.Type,
			}
		case "DRIVING_LICENSE":
			documents.DrivingLicense = &dbInitJobApplication.Document{
				URL:  cred.URL,
				Type: cred.Type,
			}
		case "RESUME":
			documents.Resume = &dbInitJobApplication.Document{
				URL:  cred.URL,
				Type: cred.Type,
			}
		}
	}

	return documents
}

func getJobApplicationDocuments(initJobApplication *dbInitJobApplication.InitJobApplication) dbJobApplication.Documents {
	var documents dbJobApplication.Documents

	if initJobApplication.ApplicantDetails.Documents.PANCard != nil {
		documents.PANCard = &dbJobApplication.Document{
			URL:  initJobApplication.ApplicantDetails.Documents.PANCard.URL,
			Type: initJobApplication.ApplicantDetails.Documents.PANCard.Type,
		}
	}

	if initJobApplication.ApplicantDetails.Documents.AadharCard != nil {
		documents.AadharCard = &dbJobApplication.Document{
			URL:  initJobApplication.ApplicantDetails.Documents.AadharCard.URL,
			Type: initJobApplication.ApplicantDetails.Documents.AadharCard.Type,
		}
	}

	if initJobApplication.ApplicantDetails.Documents.Passport != nil {
		documents.Passport = &dbJobApplication.Document{
			URL:  initJobApplication.ApplicantDetails.Documents.Passport.URL,
			Type: initJobApplication.ApplicantDetails.Documents.Passport.Type,
		}
	}

	if initJobApplication.ApplicantDetails.Documents.DrivingLicense != nil {
		documents.DrivingLicense = &dbJobApplication.Document{
			URL:  initJobApplication.ApplicantDetails.Documents.DrivingLicense.URL,
			Type: initJobApplication.ApplicantDetails.Documents.DrivingLicense.Type,
		}
	}

	if initJobApplication.ApplicantDetails.Documents.Resume != nil {
		documents.Resume = &dbJobApplication.Document{
			URL:  initJobApplication.ApplicantDetails.Documents.Resume.URL,
			Type: initJobApplication.ApplicantDetails.Documents.Resume.Type,
		}
	}

	return documents
}

func getSearchFilter(payload *searchrequest.SearchRequest) bson.D {
	var (
		role      = payload.Message.Intent.Item.Descriptor.Name
		provider  = payload.Message.Intent.Provider.Descriptor.Name
		locations = payload.Message.Intent.Provider.Locations
		tags      = payload.Message.Intent.Item.Tags
		query     = bson.D{}
	)

	if role != "" {
		query = append(query, bson.E{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: role}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "description", Value: bson.D{{Key: "$regex", Value: role}, {Key: "$options", Value: "i"}}}},
			},
		})
	}

	if provider != "" {
		query = append(query, bson.E{Key: "business.name", Value: provider})
	}

	if locations != nil {
		var location = locations[0]

		if location.City.Code != "" {
			query = append(query, bson.E{Key: "location.city", Value: location.City.Code})
		}
		if location.State.Code != "" {
			query = append(query, bson.E{Key: "location.state", Value: location.State.Code})
		}
		if location.AreaCode.Code != "" {
			query = append(query, bson.E{Key: "location.area_code", Value: location.AreaCode.Code})
			query = append(query, bson.E{Key: "$or", Value: bson.A{bson.D{{
				Key: "location.address",
				Value: bson.D{
					{Key: "$regex", Value: location.AreaCode.Code},
					{Key: "$options", Value: "i"},
				},
			}}}})
		}
		if location.Coordinates.Longitute != 0 && location.Coordinates.Latitude != 0 {
			query = append(query, bson.E{
				Key: "address.coordinates", Value: bson.D{
					{Key: "$nearSphere", Value: bson.D{
						{Key: "$geometry", Value: bson.D{
							{Key: "type", Value: "Point"},
							{
								Key: "coordinates",
								Value: bson.A{
									location.Coordinates.Longitute,
									location.Coordinates.Latitude,
								},
							},
						}},
						{Key: "$maxDistance", Value: 5000}, // 5 km in meters
					}},
				},
			})
		}
	}

	for _, tag := range tags {
		if tag.Descriptor.Code == "JOB_DETAILS" {
			for _, listItem := range tag.List {
				if listItem.Descriptor.Code == "INDUSTRY_TYPE" {
					query = append(query, bson.E{Key: "business.industry", Value: listItem.Value})
				}
				if listItem.Descriptor.Code == "JOB_TYPE" {
					query = append(query, bson.E{Key: "type", Value: listItem.Value})
				}
			}
		}
	}

	return query
}
