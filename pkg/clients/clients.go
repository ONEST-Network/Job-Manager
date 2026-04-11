package clients

import (
	apiclient "github.com/ONEST-Network/Job-Manager-Adapter/pkg/api-client"
	dbBusiness "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/business"
	dbInitJobApplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/init-job-application"
	dbJob "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job"
	dbJobApplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job-application"
	dbScheme "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/scheme"
	dbSchemeApplication "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/scheme-application"
)

type Clients struct {
	ApiClient                apiclient.Interface
	JobClient                *dbJob.Dao
	BusinessClient           *dbBusiness.Dao
	JobApplicationClient     *dbJobApplication.Dao
	InitJobApplicationClient *dbInitJobApplication.Dao
	SchemeClient             *dbScheme.Dao
	SchemeApplicationClient  *dbSchemeApplication.Dao
}

func NewClients(jobClient *dbJob.Dao, businessClient *dbBusiness.Dao, jobApplicationClient *dbJobApplication.Dao, initJobApplicationClient *dbInitJobApplication.Dao, schemeClient *dbScheme.Dao, schemeApplicationClient *dbSchemeApplication.Dao) *Clients {
	return &Clients{
		ApiClient:                apiclient.NewAPIClient(),
		JobClient:                jobClient,
		BusinessClient:           businessClient,
		JobApplicationClient:     jobApplicationClient,
		InitJobApplicationClient: initJobApplicationClient,
		SchemeClient:             schemeClient,
		SchemeApplicationClient:  schemeApplicationClient,
	}
}
