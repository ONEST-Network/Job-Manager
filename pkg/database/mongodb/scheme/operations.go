package scheme

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
	CreateScheme(scheme *Scheme) error
	GetScheme(schemeId string) (*Scheme, error)
	ListSchemes(query bson.D) ([]Scheme, error)
	UpdateScheme(query, update bson.D) error
	DeleteScheme(schemeId string) error
}

type Dao struct {
	collection *mongo.Collection
}

func NewSchemeDao(collection *mongo.Collection) *Dao {
	if err := ensure2dsphereIndex(collection, "scheme_location_2dsphere_index"); err != nil {
		logrus.Fatalf("Failed to create 2dsphere index for %s collection, %v", collection.Name(), err)
	}
	return &Dao{
		collection: collection,
	}
}

const dbTimeout = 10 * time.Second

func (d *Dao) CreateScheme(scheme *Scheme) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := database.Operator.Create(ctx, d.collection, scheme); err != nil {
		return err
	}

	return nil
}

func (d *Dao) GetScheme(schemeId string) (*Scheme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query = bson.D{{Key: "id", Value: schemeId}}

	result := database.Operator.Get(ctx, d.collection, query)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var scheme Scheme
	if err := result.Decode(&scheme); err != nil {
		return nil, err
	}

	return &scheme, nil
}

func (d *Dao) ListSchemes(query bson.D) ([]Scheme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := database.Operator.List(ctx, d.collection, query)
	if err != nil {
		return nil, err
	}

	var schemes []Scheme
	if err = result.All(ctx, &schemes); err != nil {
		return nil, err
	}

	return schemes, nil
}

func (d *Dao) UpdateScheme(query, update bson.D) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := database.Operator.Update(ctx, d.collection, query, update); err != nil {
		return err
	}

	return nil
}

func (d *Dao) DeleteScheme(schemeId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := bson.D{{Key: "id", Value: schemeId}}

	if _, err := database.Operator.Delete(ctx, d.collection, query); err != nil {
		return err
	}

	return nil
}

func ensure2dsphereIndex(collection *mongo.Collection, indexName string) error {
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

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "location.coordinates", Value: "2dsphere"}},
		Options: options.Index().SetName(indexName),
	}

	if _, err = collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		return err
	}

	logrus.Infof("2dsphere index %s created for %s collection", indexName, collection.Name())
	return nil
}
