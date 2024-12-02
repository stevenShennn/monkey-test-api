package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) *mongo.Database {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)

	return client.Database("test_db")
}

func cleanupTestDB(t *testing.T, db *mongo.Database) {
	ctx := context.Background()
	err := db.Drop(ctx)
	assert.NoError(t, err)
}

func TestMongoStore(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	store := NewMongoStore(db)
	ctx := context.Background()

	// 测试插入父请求
	req := &Request{
		RequestID:  "test-123",
		Method:     "GET",
		URL:        "http://example.com",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Timestamp:  time.Now(),
	}

	err := store.InsertRequest(ctx, req)
	assert.NoError(t, err)

	// 测试获取父请求
	fetchedReq, err := store.GetRequestByID(ctx, req.RequestID)
	assert.NoError(t, err)
	assert.Equal(t, req.RequestID, fetchedReq.RequestID)

	// 测试插入子请求
	testObjs := []TestObject{
		{
			TestID:         "test-obj-1",
			ParentRequestID: req.RequestID,
			Method:         req.Method,
			URL:            req.URL,
			Timestamp:      time.Now(),
		},
	}

	err = store.InsertTestObjects(ctx, testObjs)
	assert.NoError(t, err)

	// 测试更新子请求
	updateData := bson.M{"reason": "test update"}
	err = store.UpdateTestObjectsByRequestID(ctx, req.RequestID, updateData)
	assert.NoError(t, err)

	// 测试删除父请求（包括关联的子请求）
	err = store.DeleteRequestByID(ctx, req.RequestID)
	assert.NoError(t, err)

	// 验证删除结果
	fetchedReq, err = store.GetRequestByID(ctx, req.RequestID)
	assert.NoError(t, err)
	assert.Nil(t, fetchedReq)
} 