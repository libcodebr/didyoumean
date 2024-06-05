package mongodb

import (
	"context"
	"github.com/libcodebr/didyoumean/entity"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockMongoDB struct {
	mock.Mock
}

func (m *MockMongoDB) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMongoDB) CreateIndex(ctx context.Context, collection string, fields []string, options *options.IndexOptions) error {
	args := m.Called(ctx, collection, fields, options)
	return args.Error(0)
}

func (m *MockMongoDB) Find(ctx context.Context, collection string, fields []string, query string, pagination *entity.Pagination, result interface{}) error {
	args := m.Called(ctx, collection, fields, query, pagination, result)
	return args.Error(0)
}

func (m *MockMongoDB) UpsertBatch(ctx context.Context, collection string, documents []interface{}) (*mongo.BulkWriteResult, error) {
	args := m.Called(ctx, collection, documents)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.BulkWriteResult), args.Error(1)
}
