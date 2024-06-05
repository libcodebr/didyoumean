package mongodb

import (
	"context"
	"errors"
	"github.com/libcodebr/didyoumean/entity"
	v "github.com/libcodebr/didyoumean/pkg/verifier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type MongoDB interface {
	Close(ctx context.Context) error
	Find(ctx context.Context, collection string, fields []string, query string, pagination *entity.Pagination, result interface{}) error
	CreateIndex(ctx context.Context, collection string, fields bson.D, options *options.IndexOptions) error
	UpsertBatch(ctx context.Context, collection string, documents []interface{}) (*mongo.BulkWriteResult, error)
}

var (
	ErrDocumentDoesNotHaveIDField = errors.New("document does not have a id field")
	ErrQueryIsEmpty               = errors.New("query cannot be empty")
)

// NewMongoDB creates a new MongoDB instance.
func NewMongoDB(ctx context.Context, cfg *Config) (MongoDB, error) {
	opts := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	m := &mongoDB{
		cfg: cfg,
		db:  client.Database(cfg.Database),
	}

	if err := v.Verifier.Struct(m); err != nil {
		return nil, err
	}

	return m, nil
}

// mongoDB is a MongoDB implementation.
type mongoDB struct {
	db  *mongo.Database `validate:"required"`
	cfg *Config         `validate:"required"`
}

// Close closes the MongoDB connection.
func (m *mongoDB) Close(ctx context.Context) error {
	return m.db.Client().Disconnect(ctx)
}

// CreateIndex creates an index for the given collection and field.
func (m *mongoDB) CreateIndex(ctx context.Context, collection string, fields bson.D, options *options.IndexOptions) error {
	idxModel := mongo.IndexModel{
		Keys:    fields,
		Options: options,
	}

	_, err := m.db.Collection(collection).Indexes().CreateOne(ctx, idxModel)
	return err
}

// Find finds documents in a collection with query and pagination.
func (m *mongoDB) Find(ctx context.Context, collection string, fields []string, query string, pagination *entity.Pagination, result interface{}) error {
	if query == "" || strings.TrimSpace(query) == "" {
		return ErrQueryIsEmpty
	}

	sf := make([]bson.M, 0, len(fields))
	for _, field := range fields {
		sf = append(sf, bson.M{field: bson.M{"$regex": query, "$options": "i"}})
	}
	filter := bson.M{"$or": sf}

	total, err := m.db.Collection(collection).CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	pagination.TotalDocuments = total

	offset := (pagination.Page - 1) * pagination.PageSize
	if offset >= pagination.TotalDocuments {
		pagination.Page = 1
		offset = 0
	}

	opts := options.Find().SetSkip(offset).SetLimit(pagination.PageSize)
	cursor, err := m.db.Collection(collection).Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}

// UpsertBatch upserts a batch of documents into a collection.
func (m *mongoDB) UpsertBatch(ctx context.Context, collection string, documents []interface{}) (*mongo.BulkWriteResult, error) {
	coll := m.db.Collection(collection)
	models := make([]mongo.WriteModel, len(documents))

	for i, doc := range documents {
		bsonDoc, err := bson.Marshal(doc)
		if err != nil {
			return nil, err
		}

		var docMap primitive.M
		err = bson.Unmarshal(bsonDoc, &docMap)
		if err != nil {
			return nil, err
		}

		id, ok := docMap["id"]
		if !ok {
			return nil, ErrDocumentDoesNotHaveIDField
		}

		filter := bson.M{"id": id}

		models[i] = mongo.NewReplaceOneModel().SetFilter(filter).SetReplacement(docMap).SetUpsert(true)
	}

	opts := options.BulkWrite().SetOrdered(false)
	result, err := coll.BulkWrite(ctx, models, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}
