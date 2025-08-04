package jobapplication

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	database "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb"
)

type DaoInterface interface {
	CreateJobApplication(jobApplication *JobApplication) error
	GetJobApplication(jobApplicationID string) (*JobApplication, error)
	ListJobApplication(query bson.D) ([]JobApplication, error)
	DeleteJobApplication(applicationID, name string) error
	UpdateJobApplication(query, update bson.D) error
}

type Dao struct {
	collection *mongo.Collection
}

func NewJobApplicationDao(collection *mongo.Collection) *Dao {
	return &Dao{
		collection: collection,
	}
}

const dbTimeout = 10 * time.Second

// CreateJobApplication creates a job application post in the database
func (d *Dao) CreateJobApplication(jobApplication *JobApplication) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := database.Operator.Create(ctx, d.collection, jobApplication); err != nil {
		return err
	}

	return nil
}

func (d *Dao) GetJobApplication(jobApplicationID string) (*JobApplication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query = bson.D{{Key: "id", Value: jobApplicationID}}

	result := database.Operator.Get(ctx, d.collection, query)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var jobApplication JobApplication
	if err := result.Decode(&jobApplication); err != nil {
		return nil, err
	}

	return &jobApplication, nil
}

func (d *Dao) ListJobApplication(query bson.D) ([]JobApplication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	cursor, err := database.Operator.List(ctx, d.collection, query)
	if err != nil {
		return nil, err
	}

	var jobApplications []JobApplication
	if err := cursor.All(ctx, &jobApplications); err != nil {
		return nil, err
	}

	return jobApplications, nil
}

func (d *Dao) DeleteJobApplication(applicationID, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query = bson.D{{Key: "id", Value: applicationID}}

	if _, err := database.Operator.Delete(ctx, d.collection, query); err != nil {
		return err
	}

	return nil
}

func (d *Dao) UpdateJobApplication(query, update bson.D) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := database.Operator.Update(ctx, d.collection, query, update); err != nil {
		return err
	}

	return nil
}

func (d *Dao) UpdateJobApplicationAndReturnDocument(query, update bson.D) (*JobApplication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var jobApplication JobApplication
	if err := database.Operator.UpdateAndReturnDocument(ctx, d.collection, query, update).Decode(&jobApplication); err != nil {
		return nil, err
	}

	return &jobApplication, nil
}
