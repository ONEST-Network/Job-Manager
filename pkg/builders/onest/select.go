package onest

import (
	"time"

	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/select/request"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/select/response"
)

func BuildSendJobFulfillmentResponse(payload *request.SelectRequest) *response.SelectResponse {
	res := response.SelectResponse{
		Context: response.Context{
			Domain:        payload.Context.Domain,
			Action:        "on_select",
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
				Provider: response.Provider{
					ID: payload.Message.Order.Provider.ID,
				},
				Fulfillments: []response.Fulfillments{
					{
						ID:   "F1",
						Type: "lead & recruitment",
					},
				},
				Items: []response.Items{},
			},
		},
	}

	for _, item := range payload.Message.Order.Items {
		res.Message.Order.Items = append(res.Message.Order.Items, response.Items{
			ID:             item.ID,
			FulfillmentIds: []string{"F1"},
		})
	}

	return &res
}
