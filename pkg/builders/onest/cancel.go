package onest

import (
	"time"

	dbJobApplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job-application"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/cancel/request"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/cancel/response"
)

func BuildWithdrawJobApplicationResponse(payload *request.CancelRequest, jobApplication *dbJobApplication.JobApplication) *response.CancelResponse {
	return &response.CancelResponse{
		Context: response.Context{
			Domain:        payload.Context.Domain,
			Action:        "on_cancel",
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
				ID:     payload.Message.OrderID,
				Status: "Cancelled",
				Provider: response.Provider{
					ID: "1",
				},
				Items: []response.Items{
					{
						ID:             jobApplication.JobID,
						FulfillmentIds: []string{"F1", "C1"},
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
								Code: "CANCELLED",
							},
						},
					},
					{
						ID: "C1",
					},
				},
			},
		},
	}
}
