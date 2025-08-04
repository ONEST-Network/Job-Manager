package job

import "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/business"

// Job represents a job in the database
type Job struct {
	ID             string            `bson:"id"`
	Name           string            `bson:"name" json:"name"`
	Description    string            `bson:"description" json:"description"`
	Type           JobType           `bson:"type" json:"type"`
	Vacancies      int               `bson:"vacancies" json:"vacancies"`
	SalaryRange    SalaryRange       `bson:"salary_range" json:"salaryRange"`
	ApplicationIDs []string          `bson:"application_ids" json:"applicationIds"`
	Business       business.Business `bson:"business" json:"business"`
	WorkHours      WorkHours         `bson:"work_hours" json:"workHours"`
	WorkDays       WorkDays          `bson:"work_days" json:"workDays"`
	Eligibility    Eligibility       `bson:"eligibility" json:"eligibility"`
	Location       Location          `bson:"location" json:"location"`
}

// SalaryRange represents the salary range of a job
type SalaryRange struct {
	Min int `bson:"min" json:"min"`
	Max int `bson:"max" json:"max"`
}

type JobType string

const (
	JobTypeFullTime   JobType = "full-time"
	JobTypePartTime   JobType = "part-time"
	JobTypeContract   JobType = "contract"
	JobTypeInternship JobType = "internship"
)

type Eligibility struct {
	Gender                Gender                `bson:"gender" json:"gender"`
	YearsOfExperience     int                   `bson:"years_of_experience" json:"yearsOfExperience"`
	DocumentsRequired     []Document            `bson:"documents_required" json:"documentsRequired"`
	AcademicQualification AcademicQualification `bson:"academic_qualification" json:"academicQualification"`
}

type Gender string

const (
	GenderAny    Gender = "any"
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

// Document represents the document required for a job
type Document string

const (
	DocumentAadharCard            Document = "aadhar_card"
	DocumentPanCard               Document = "pan_card"
	DocumentDrivingLic            Document = "driving_license"
	DocumentClassXCert            Document = "class_x_certificate"
	DocumentClassXIICertificate   Document = "class_xii_certificate"
	DocumentDiplomaCertificate    Document = "diploma_certificate"
	DocumentGraduationCertificate Document = "graduation_certificate"
	DocumentPostGradCertificate   Document = "post_grad_certificate"
	DocumentPassport              Document = "passport"
	DocumentOther                 Document = "other"
)

// AcademicQualification represents the academic qualification of a job
type AcademicQualification string

const (
	AcademicQualificationNone         AcademicQualification = "None"
	AcademicQualificationClassX       AcademicQualification = "Class-X"
	AcademicQualificationClassXII     AcademicQualification = "Class-XII"
	AcademicQualificationDiploma      AcademicQualification = "Diploma"
	AcademicQualificationGraduate     AcademicQualification = "Graduate"
	AcademicQualificationPostGraduate AcademicQualification = "Post-Graduate"
)

// WorkHours represents the start and end time of a job
// stored in military time format, for eg. 0900, 1800
type WorkHours struct {
	Start string `bson:"start" json:"start"`
	End   string `bson:"end" json:"end"`
}

// WorkDays represents the start and end day of a job
// 1-6, 1 is Monday, 6 is Saturday
type WorkDays struct {
	Start int `bson:"start" json:"start"`
	End   int `bson:"end" json:"end"`
}

// Location represents the location of a job
type Location struct {
	Coordinates Coordinates `bson:"coordinates" json:"coordinates"`
	Address     string      `bson:"address" json:"address"`
	Street      string      `bson:"street" json:"street"`
	AreaCode    string      `bson:"area_code" json:"areaCode"` // Postal code, for example: '560102'
	City        string      `bson:"city" json:"city"`          // STD code, for example: 'std:080'
	State       string      `bson:"state" json:"state"`        // State code, for example: 'IN-KA'
}

// Coordinates represents the longitude and latitude of a location
type Coordinates struct {
	Type        string    `bson:"type" json:"type"`               // GeoJSON type, it shall be equal to 'Point'
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // longitude, latitude
}
