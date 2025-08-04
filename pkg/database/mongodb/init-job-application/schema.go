package jobapplication

import "time"

type InitJobApplication struct {
	TransactionID    string           `bson:"transaction_id"`
	JobID            string           `bson:"job_id"`
	ApplicantDetails ApplicantDetails `bson:"applicant_details"`
	CreatedAt        time.Time        `bson:"created_at"`
}

type ApplicantDetails struct {
	Name       string     `bson:"name"`
	Gender     string     `bson:"gender"`
	Age        int        `bson:"age"`
	Experience Experience `bson:"experience"`
	Documents  Documents  `bson:"documents"`
	Phone      string     `bson:"phone"`
	Email      string     `bson:"email"`
}

type Documents struct {
	PANCard        *Document `bson:"pan_card"`
	AadharCard     *Document `bson:"aadhar_card"`
	Passport       *Document `bson:"passport"`
	DrivingLicense *Document `bson:"driving_license"`
	Resume         *Document `bson:"resume"`
}

type Document struct {
	URL  string `bson:"url"`
	Type string `bson:"type"`
}

type Experience struct {
	Years int `bson:"years"`
}
