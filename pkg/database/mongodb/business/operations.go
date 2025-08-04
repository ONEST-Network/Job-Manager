package business

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DaoInterface interface {
	GetBusiness(id string) (*Business, error)
	CreateBusiness(business *Business) error
	ListBusinesses(query bson.D) ([]Business, error)
}

type Dao struct {
	collection *mongo.Collection
}

func NewBusinessDao(collection *mongo.Collection) *Dao {
	return &Dao{
		collection: collection,
	}
}

const dbTimeout = 10 * time.Second

func (d *Dao) GetBusiness(id string) (*Business, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var business Business
	if err := d.collection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&business); err != nil {
		return nil, err
	}

	return &business, nil
}

func (d *Dao) CreateBusiness(business *Business) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := d.collection.InsertOne(ctx, business); err != nil {
		return err
	}

	return nil
}

func (d *Dao) ListBusinesses(query bson.D) ([]Business, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	cursor, err := d.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var businesses []Business
	if err := cursor.All(ctx, &businesses); err != nil {
		return nil, err
	}

	return businesses, nil
}
