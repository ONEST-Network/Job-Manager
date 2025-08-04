package job

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	database "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb"
	"github.com/sirupsen/logrus"
)

type DaoInterface interface {
	CreateJob(job *Job) error
	GetJob(jobId string) error
	ListJobs(query bson.D) ([]Job, error)
	DeleteJob(jobID string) error
	UpdateJob(query, update bson.D) error
}

type Dao struct {
	collection *mongo.Collection
}

func NewJobDao(collection *mongo.Collection) *Dao {
	if err := ensure2dsphereIndex(collection, "coordinates_2dsphere_index"); err != nil {
		logrus.Fatalf("Failed to create 2dsphere index for %s collection, %v", collection.Name(), err)
	}
	return &Dao{
		collection: collection,
	}
}

const dbTimeout = 10 * time.Second

// CreateJob creates a job post in the database
func (d *Dao) CreateJob(job *Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := database.Operator.Create(ctx, d.collection, job); err != nil {
		return err
	}

	return nil
}

// GetJob returns the job for the given id
func (d *Dao) GetJob(jobId string) (*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query = bson.D{{Key: "id", Value: jobId}}

	result := database.Operator.Get(ctx, d.collection, query)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var job Job
	if err := result.Decode(&job); err != nil {
		return nil, err
	}

	return &job, nil
}

// ListJobs lists jobs from the database
func (d *Dao) ListJobs(query bson.D) ([]Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := database.Operator.List(ctx, d.collection, query)
	if err != nil {
		return nil, err
	}

	var jobs []Job
	if err = result.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

// UpdateJob updates a job in the database
func (d *Dao) UpdateJob(query, update bson.D) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := database.Operator.Update(ctx, d.collection, query, update); err != nil {
		return err
	}

	return nil
}

// DeleteJob deletes a job from the database
func (d *Dao) DeleteJob(jobID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := bson.D{{Key: "id", Value: jobID}}

	if _, err := database.Operator.Delete(ctx, d.collection, query); err != nil {
		return err
	}

	return nil
}

func ensure2dsphereIndex(collection *mongo.Collection, indexName string) error {
	ctx := context.Background()

	// Check if the index already exists
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			return err
		}
		if name, ok := index["name"].(string); ok && name == indexName {
			// index already exists
			return nil
		}
	}

	// Create the 2dsphere index
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "location.coordinates", Value: "2dsphere"}}, // Index on location.coordinates field
		Options: options.Index().SetName(indexName),
	}

	if _, err = collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		return err
	}

	logrus.Infof("2dsphere index %s created for %s collection", indexName, collection.Name())
	return nil
}
