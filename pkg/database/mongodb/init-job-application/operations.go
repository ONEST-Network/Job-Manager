package jobapplication

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	database "github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb"
)

type DaoInterface interface {
	CreateInitJobApplication(jobApplication *InitJobApplication) error
	GetInitJobApplication(query bson.D) (*InitJobApplication, error)
	ListInitJobApplication(query bson.D) ([]InitJobApplication, error)
	DeleteInitJobApplication(applicationID, name string) error
	UpdateInitJobApplication(query, update bson.D) error
}

type Dao struct {
	collection *mongo.Collection
}

const dbTimeout = 10 * time.Second

func NewInitJobApplicationDao(collection *mongo.Collection) *Dao {
	if err := ensureTTLIndex(collection, "created_at_ttl_index", 3600); err != nil {
		logrus.Fatalf("Failed to create TTL index for %s collection, %v", collection.Name(), err)
	}
	return &Dao{
		collection: collection,
	}
}

func (d *Dao) CreateInitJobApplication(jobApplication *InitJobApplication) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := database.Operator.Create(ctx, d.collection, jobApplication); err != nil {
		return err
	}

	return nil
}

func (d *Dao) GetInitJobApplication(transactionId string) (*InitJobApplication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query = bson.D{{Key: "transaction_id", Value: transactionId}}

	result := database.Operator.Get(ctx, d.collection, query)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var initJobApplication InitJobApplication
	if err := result.Decode(&initJobApplication); err != nil {
		return nil, err
	}

	return &initJobApplication, nil
}

func (d *Dao) ListInitJobApplication(query bson.D) ([]InitJobApplication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	cursor, err := database.Operator.List(ctx, d.collection, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var initJobApplications []InitJobApplication
	if err := cursor.All(ctx, &initJobApplications); err != nil {
		return nil, err
	}

	return initJobApplications, nil
}

func (d *Dao) DeleteInitJobApplication(transactionId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query = bson.D{{Key: "transaction_id", Value: transactionId}}

	if _, err := database.Operator.Delete(ctx, d.collection, query); err != nil {
		return err
	}

	return nil
}

func ensureTTLIndex(collection *mongo.Collection, indexName string, expireSeconds int32) error {
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
			// TTL index already exists
			return nil
		}
	}

	// Create the TTL index
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: 1}}, // Index on created_at field
		Options: options.Index().SetExpireAfterSeconds(expireSeconds).SetName(indexName),
	}

	if _, err = collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		return err
	}

	logrus.Infof("TTL index %s created for %s collection", indexName, collection.Name())
	return nil
}
