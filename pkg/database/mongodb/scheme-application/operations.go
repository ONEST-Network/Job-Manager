package schemeapplication

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
	CreateSchemeApplication(app *SchemeApplication) error
	GetSchemeApplication(appId string) (*SchemeApplication, error)
	ListSchemeApplications(query bson.D) ([]SchemeApplication, error)
	UpdateSchemeApplication(query, update bson.D) error
}

type Dao struct {
	collection *mongo.Collection
}

func NewSchemeApplicationDao(collection *mongo.Collection) *Dao {
	if err := ensureTTLIndex(collection, "created_at_ttl_index"); err != nil {
		logrus.Fatalf("Failed to create TTL index for %s collection, %v", collection.Name(), err)
	}
	return &Dao{
		collection: collection,
	}
}

const dbTimeout = 10 * time.Second

func (d *Dao) CreateSchemeApplication(app *SchemeApplication) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := database.Operator.Create(ctx, d.collection, app); err != nil {
		return err
	}

	return nil
}

func (d *Dao) GetSchemeApplication(appId string) (*SchemeApplication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query = bson.D{{Key: "id", Value: appId}}

	result := database.Operator.Get(ctx, d.collection, query)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var app SchemeApplication
	if err := result.Decode(&app); err != nil {
		return nil, err
	}

	return &app, nil
}

func (d *Dao) ListSchemeApplications(query bson.D) ([]SchemeApplication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := database.Operator.List(ctx, d.collection, query)
	if err != nil {
		return nil, err
	}

	var apps []SchemeApplication
	if err = result.All(ctx, &apps); err != nil {
		return nil, err
	}

	return apps, nil
}

func (d *Dao) UpdateSchemeApplication(query, update bson.D) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := database.Operator.Update(ctx, d.collection, query, update); err != nil {
		return err
	}

	return nil
}

func ensureTTLIndex(collection *mongo.Collection, indexName string) error {
	ctx := context.Background()

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
			return nil
		}
	}

	// Index will expire documents after 30 days
	expireAfterSeconds := int32(30 * 24 * 60 * 60)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().SetName(indexName).SetExpireAfterSeconds(expireAfterSeconds),
	}

	if _, err = collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		return err
	}

	logrus.Infof("TTL index %s created for %s collection", indexName, collection.Name())
	return nil
}
