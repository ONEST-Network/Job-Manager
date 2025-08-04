package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoOperator interface {
	Create(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error)
	Get(ctx context.Context, collection *mongo.Collection, query bson.D) *mongo.SingleResult
	List(ctx context.Context, collection *mongo.Collection, query bson.D) (*mongo.Cursor, error)
	Update(ctx context.Context, collection *mongo.Collection, query, update bson.D,
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, collection *mongo.Collection, query, update bson.D,
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateAndReturnDocument(ctx context.Context, collection *mongo.Collection, query, update bson.D) *mongo.SingleResult
	Delete(ctx context.Context, collection *mongo.Collection, query bson.D, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Aggregate(ctx context.Context, collection *mongo.Collection, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	ListDataBase(ctx context.Context, mclient *mongo.Client) ([]string, error)
}

type MongoOperations struct{}

var (
	// Operator contains all the CRUD operations of the mongo database
	Operator MongoOperator = &MongoOperations{}
)

// Create puts a document in the database
func (m *MongoOperations) Create(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	return collection.InsertOne(ctx, document)
}

// Get fetches a document from the database based on a query
func (m *MongoOperations) Get(ctx context.Context, collection *mongo.Collection, query bson.D) *mongo.SingleResult {
	return collection.FindOne(ctx, query)
}

// List fetches a list of documents from the database based on a query
func (m *MongoOperations) List(ctx context.Context, collection *mongo.Collection, query bson.D) (*mongo.Cursor, error) {
	return collection.Find(ctx, query)
}

// Update updates a document in the database based on a query
func (m *MongoOperations) Update(ctx context.Context, collection *mongo.Collection, query, update bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return collection.UpdateOne(ctx, query, update, opts...)
}

// UpdateAndReturnDocument updates a document and then returns the updated document
func (m *MongoOperations) UpdateAndReturnDocument(ctx context.Context, collection *mongo.Collection, query, update bson.D) *mongo.SingleResult {
	return collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
}

// Update updates a document in the database based on a query
func (m *MongoOperations) UpdateMany(ctx context.Context, collection *mongo.Collection, query, update bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return collection.UpdateMany(ctx, query, update, opts...)
}

// Delete removes a document from the database based on a query
func (m *MongoOperations) Delete(ctx context.Context, collection *mongo.Collection, query bson.D, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return collection.DeleteOne(ctx, query, opts...)
}

func (m *MongoOperations) Aggregate(ctx context.Context, collection *mongo.Collection, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	result, err := collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MongoOperations) ListDataBase(ctx context.Context, mclient *mongo.Client) ([]string, error) {
	dbs, err := mclient.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	return dbs, nil
}
