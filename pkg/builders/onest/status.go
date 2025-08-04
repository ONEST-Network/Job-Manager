package onest

import (
	"time"

	dbJobApplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job-application"

	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/status/request"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/status/response"
)

func BuildJobApplicationStatusResponse(payload *request.StatusRequest, jobApplication *dbJobApplication.JobApplication) *response.StatusResponse {
	return &response.StatusResponse{
		Context: response.Context{
			Domain:        payload.Context.Domain,
			Action:        "on_status",
			Version:       payload.Context.Version,
			BapID:         payload.Context.BapID,
			BapURI:        payload.Context.BapURI,
			BppID:         payload.Context.BppID,
			BppURI:        payload.Context.BppURI,
			TransactionID: payload.Context.TransactionID,
			MessageID:     payload.Context.MessageID,
			Location: response.Location{
				City: response.City{
					Code: payload.Context.Location.City.Code,
				},
				Country: response.Country{
					Code: payload.Context.Location.Country.Code,
				},
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			TTL:       "PT30S",
		},
		Message: response.Message{
			Order: response.Order{
				ID:     jobApplication.ID,
				Status: string(jobApplication.Status),
				Provider: response.Provider{
					ID: "1",
				},
				Items: []response.Items{
					{
						ID:             jobApplication.JobID,
						FulfillmentIds: []string{"F1"},
						Time: response.Time{
							Range: response.Range{
								Start: jobApplication.CreatedAt.UTC().Format(time.RFC3339),
								End:   jobApplication.CreatedAt.Add(time.Hour * 24 * 30).UTC().Format(time.RFC3339),
							},
						},
					},
				},
				Fulfillments: []response.Fulfillments{
					{
						ID:   "F1",
						Type: "lead & recruitment",
						State: response.State{
							Descriptor: response.Descriptor{
								Code: string(jobApplication.Status),
							},
							UpdatedAt: jobApplication.UpdatedAt.UTC().Format(time.RFC3339),
						},
					},
				},
			},
		},
	}
}
