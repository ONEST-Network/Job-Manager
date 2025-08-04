package jobapplication

type UpdateJobApplicationStatusRequest struct {
	// @Enum(APPLICATION_ACCEPTED, APPLICATION_REJECTED, ASSESSMENT_IN_PROGRESS, OFFER_REJECTED, OFFER_ACCEPTED, OFFER_EXTENDED, CANCELLED)
	Status string `json:"status"`
}
