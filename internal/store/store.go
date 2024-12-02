package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store 定义存储层接口
type Store interface {
	// 父请求操作
	InsertRequest(ctx context.Context, req *Request) error
	GetRequestByID(ctx context.Context, requestID string) (*Request, error)
	DeleteRequestByID(ctx context.Context, requestID string) error
	GetRequestsByTime(ctx context.Context, limit int64) ([]Request, error)

	// 子请求操作
	InsertTestObjects(ctx context.Context, testObjs []TestObject) error
	DeleteTestObjectsByRequestID(ctx context.Context, requestID string) error
	UpdateTestObjectsByRequestID(ctx context.Context, requestID string, updateData bson.M) error
}

// MongoStore MongoDB 存储实现
type MongoStore struct {
	db *mongo.Database
}

// NewMongoStore 创建新的 MongoDB 存储实例
func NewMongoStore(db *mongo.Database) Store {
	return &MongoStore{db: db}
} 