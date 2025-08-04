package onest

import (
	"time"

	confirmrequest "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/confirm/request"
	confirmresponse "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/confirm/response"
)

func BuildConfirmJobApplicationResponse(payload *confirmrequest.ConfirmRequest) *confirmresponse.ConfirmResponse {
	res := confirmresponse.ConfirmResponse{
		Context: confirmresponse.Context{
			Domain:        payload.Context.Domain,
			Action:        "on_confirm",
			Version:       payload.Context.Version,
			BapID:         payload.Context.BapID,
			BapURI:        payload.Context.BapURI,
			BppID:         payload.Context.BppID,
			BppURI:        payload.Context.BppURI,
			TransactionID: payload.Context.TransactionID,
			MessageID:     payload.Context.MessageID,
			Location: confirmresponse.Location{
				City: confirmresponse.City{
					Code: payload.Context.Location.City.Code,
				},
				Country: confirmresponse.Country{
					Code: payload.Context.Location.Country.Code,
				},
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			TTL:       "PT30S",
		},
		Message: confirmresponse.Message{
			Order: confirmresponse.Order{
				ID: payload.Message.Order.ID,
				Provider: confirmresponse.Provider{
					ID: payload.Message.Order.Provider.ID,
				},
				Items: getConfirmItems(payload),
				Fulfillments: []confirmresponse.Fulfillments{
					{
						ID:   "F1",
						Type: "lead & recruitment",
						State: confirmresponse.State{
							Descriptor: confirmresponse.StateDescriptor{
								Code: "APPLICATION_ACCEPTED",
							},
							UpdatedAt: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						},
						Customer: confirmresponse.Customer{
							Person: confirmresponse.Person{
								Name:      payload.Message.Order.Fulfillments[0].Customer.Person.Name,
								Gender:    payload.Message.Order.Fulfillments[0].Customer.Person.Gender,
								Age:       payload.Message.Order.Fulfillments[0].Customer.Person.Age,
								Skills:    getConfirmSkills(payload),
								Languages: getConfirmLanguages(payload),
								Creds:     getConfirmCreds(payload),
							},
							Contact: confirmresponse.Contact{
								Phone: payload.Message.Order.Fulfillments[0].Customer.Contact.Phone,
								Email: payload.Message.Order.Fulfillments[0].Customer.Contact.Email,
							},
						},
					},
				},
			},
		},
	}
	return &res
}

func getConfirmItems(payload *confirmrequest.ConfirmRequest) []confirmresponse.Items {
	var items []confirmresponse.Items

	for _, item := range payload.Message.Order.Items {
		items = append(items, confirmresponse.Items{
			ID:             item.ID,
			FulfillmentIds: []string{"F1"},
		})
	}

	return items
}

func getConfirmSkills(payload *confirmrequest.ConfirmRequest) []confirmresponse.Skills {
	var skills []confirmresponse.Skills

	for _, skill := range payload.Message.Order.Fulfillments[0].Customer.Person.Skills {
		skills = append(skills, confirmresponse.Skills{
			Name: skill.Name,
		})
	}

	return skills
}

func getConfirmLanguages(payload *confirmrequest.ConfirmRequest) []confirmresponse.Languages {
	var languages []confirmresponse.Languages

	for _, language := range payload.Message.Order.Fulfillments[0].Customer.Person.Languages {
		languages = append(languages, confirmresponse.Languages{
			Name: language.Name,
		})
	}

	return languages
}

func getConfirmCreds(payload *confirmrequest.ConfirmRequest) []confirmresponse.Creds {
	var creds []confirmresponse.Creds

	for _, cred := range payload.Message.Order.Fulfillments[0].Customer.Person.Creds {
		creds = append(creds, confirmresponse.Creds{
			ID: cred.ID,
			Descriptor: confirmresponse.CredsDescriptor{
				Name:      cred.Descriptor.Name,
				ShortDesc: cred.Descriptor.ShortDesc,
				LongDesc:  cred.Descriptor.LongDesc,
			},
			URL:  cred.URL,
			Type: cred.Type,
		})
	}

	return creds
}
