package business

// Business represents a business in the database
type Business struct {
	ID             string   `bson:"id"`
	Name           string   `bson:"name"`
	Phone          string   `bson:"phone"`
	Email          string   `bson:"email"`
	PictureURLs    []string `bson:"picture_urls"`
	Description    string   `bson:"description"`
	GSTIndexNumber string   `bson:"gst_index_number"`
	Location       Location `bson:"location"`
	Industry       Industry `bson:"industry"`
}

// Industry represents the industry of a business
type Industry string

const (
	IndustryRetailAndEcommerce           Industry = "RetailAndEcommerce"
	IndustryFoodAndBeverages             Industry = "FoodAndBeverages"
	IndustryHealthAndWellness            Industry = "HealthAndWellness"
	IndustryEducationAndTraining         Industry = "EducationAndTraining"
	IndustryProfessionalServices         Industry = "ProfessionalServices"
	IndustryManufacturing                Industry = "Manufacturing"
	IndustryHospitalityAndTourism        Industry = "HospitalityAndTourism"
	IndustryArtsAndEntertainment         Industry = "ArtsAndEntertainment"
	IndustryTechnologyAndSoftware        Industry = "TechnologyAndSoftware"
	IndustryConstructionAndRealEstate    Industry = "ConstructionAndRealEstate"
	IndustryTransportationAndLogistics   Industry = "TransportationAndLogistics"
	IndustryAgricultureAndFarming        Industry = "AgricultureAndFarming"
	IndustryFinanceAndInsurance          Industry = "FinanceAndInsurance"
	IndustryEnergyAndUtilities           Industry = "EnergyAndUtilities"
	IndustryNonProfitAndSocialEnterprise Industry = "NonProfitAndSocialEnterprise"
	IndustryMediaAndPublishing           Industry = "MediaAndPublishing"
	IndustryAutomotive                   Industry = "Automotive"
	IndustryFashionAndLifestyle          Industry = "FashionAndLifestyle"
	IndustrySportsAndRecreation          Industry = "SportsAndRecreation"
	IndustryOther                        Industry = "Other"
)

// Location represents the location of a business
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
