package onest

import (
	"time"

	initrequest "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/init/request"
	initresponse "github.com/ONEST-Network/Job-Manager-Adapter/pkg/types/payload/onest/init/response"
)

func BuildInitializeJobApplicationResponse(payload *initrequest.InitRequest) *initresponse.InitResponse {
	res := initresponse.InitResponse{
		Context: initresponse.Context{
			Domain:        payload.Context.Domain,
			Action:        "on_init",
			Version:       payload.Context.Version,
			BapID:         payload.Context.BapID,
			BapURI:        payload.Context.BapURI,
			BppID:         payload.Context.BppID,
			BppURI:        payload.Context.BppURI,
			TransactionID: payload.Context.TransactionID,
			MessageID:     payload.Context.MessageID,
			Location: initresponse.Location{
				City: initresponse.City{
					Code: payload.Context.Location.City.Code,
				},
				Country: initresponse.Country{
					Code: payload.Context.Location.Country.Code,
				},
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			TTL:       "PT30S",
		},
		Message: initresponse.Message{
			Order: initresponse.Order{
				Provider: initresponse.Provider{
					ID: payload.Message.Order.Provider.ID,
				},
				Items: getInitItems(payload),
				Fulfillments: []initresponse.Fulfillments{
					{
						ID:   "F1",
						Type: "lead & recruitment",
						State: initresponse.State{
							Descriptor: initresponse.StateDescriptor{
								Code: "APPLICATION_IN_PROGRESS",
							},
							UpdatedAt: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						},
						Customer: initresponse.Customer{
							Person: initresponse.Person{
								Name:      payload.Message.Order.Fulfillments[0].Customer.Person.Name,
								Gender:    payload.Message.Order.Fulfillments[0].Customer.Person.Gender,
								Age:       payload.Message.Order.Fulfillments[0].Customer.Person.Age,
								Skills:    getInitSkills(payload),
								Languages: getInitLanguages(payload),
								Creds:     getInitCreds(payload),
							},
							Contact: initresponse.Contact{
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

func getInitItems(payload *initrequest.InitRequest) []initresponse.Items {
	var items []initresponse.Items

	for _, item := range payload.Message.Order.Items {
		items = append(items, initresponse.Items{
			ID:             item.ID,
			FulfillmentIds: []string{"F1"},
		})
	}

	return items
}

func getInitSkills(payload *initrequest.InitRequest) []initresponse.Skills {
	var skills []initresponse.Skills

	for _, skill := range payload.Message.Order.Fulfillments[0].Customer.Person.Skills {
		skills = append(skills, initresponse.Skills{
			Name: skill.Name,
		})
	}

	return skills
}

func getInitLanguages(payload *initrequest.InitRequest) []initresponse.Languages {
	var languages []initresponse.Languages

	for _, language := range payload.Message.Order.Fulfillments[0].Customer.Person.Languages {
		languages = append(languages, initresponse.Languages{
			Name: language.Name,
		})
	}

	return languages
}

func getInitCreds(payload *initrequest.InitRequest) []initresponse.Creds {
	var creds []initresponse.Creds

	for _, cred := range payload.Message.Order.Fulfillments[0].Customer.Person.Creds {
		creds = append(creds, initresponse.Creds{
			ID: cred.ID,
			Descriptor: initresponse.CredsDescriptor{
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
