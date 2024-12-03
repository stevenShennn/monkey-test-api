package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"monkey-test-api/internal/logger"
)

const (
	requestCollection     = "requests"
	testObjectCollection = "test_objects"
)

// MongoStore MongoDB 存储实现
type MongoStore struct {
	db     *mongo.Database
	client *mongo.Client
}

// InsertRequest 存储父请求
func (s *MongoStore) InsertRequest(ctx context.Context, req *Request) error {
	logger.Infof("开始存储请求: %s", req.RequestID)
	_, err := s.db.Collection(requestCollection).InsertOne(ctx, req)
	if err != nil {
		logger.Errorf("存储请求失败: %v", err)
		return err
	}
	logger.Infof("请求存储成功: %s", req.RequestID)
	return nil
}

// GetRequestByID 根据 ID 获取父请求
func (s *MongoStore) GetRequestByID(ctx context.Context, requestID string) (*Request, error) {
	logger.Debugf("开始查询请求: %s", requestID)
	
	var req Request
	err := s.db.Collection(requestCollection).FindOne(ctx, bson.M{"request_id": requestID}).Decode(&req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Infof("未找到请求: %s", requestID)
			return nil, nil
		}
		logger.Errorf("查询请求失败: %v", err)
		return nil, err
	}
	
	logger.Debugf("成功查询到请求: %s", requestID)
	return &req, nil
}

// DeleteRequestByID 删除父请求及其关联的子请求
func (s *MongoStore) DeleteRequestByID(ctx context.Context, requestID string) error {
	logger.Infof("开始删除请求及其关联数据: %s", requestID)
	
	session, err := s.db.Client().StartSession()
	if err != nil {
		logger.Errorf("创建会话失败: %v", err)
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 删除父请求
		if err := s.db.Collection(requestCollection).FindOneAndDelete(sessCtx, bson.M{"request_id": requestID}).Err(); err != nil {
			logger.Errorf("删除父请求失败: %v", err)
			return nil, err
		}

		// 删除关联的子请求
		if err := s.DeleteTestObjectsByRequestID(sessCtx, requestID); err != nil {
			logger.Errorf("删除关联子请求失败: %v", err)
			return nil, err
		}

		logger.Infof("成功删除请求及其关联数据: %s", requestID)
		return nil, nil
	})

	return err
}

// GetRequestsByTime 按时间倒序获取父请求
func (s *MongoStore) GetRequestsByTime(ctx context.Context, limit int64) ([]Request, error) {
	logger.Debugf("开始查询请求列表, limit: %d", limit)
	
	opts := options.Find().
		SetSort(bson.M{"timestamp": -1}).
		SetLimit(limit)

	cursor, err := s.db.Collection(requestCollection).Find(ctx, bson.M{}, opts)
	if err != nil {
		logger.Errorf("查询请求列表失败: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []Request
	if err := cursor.All(ctx, &requests); err != nil {
		logger.Errorf("解析请求列表失败: %v", err)
		return nil, err
	}

	logger.Debugf("成功查询到 %d 条请求记录", len(requests))
	return requests, nil
}

// InsertTestObjects 批量插入子请求
func (s *MongoStore) InsertTestObjects(ctx context.Context, testObjs []TestObject) error {
	if len(testObjs) == 0 {
		logger.Debug("没有子请求需要插入")
		return nil
	}

	logger.Infof("开始批量插入 %d 个子请求", len(testObjs))
	docs := make([]interface{}, len(testObjs))
	for i, obj := range testObjs {
		docs[i] = obj
	}

	_, err := s.db.Collection(testObjectCollection).InsertMany(ctx, docs)
	if err != nil {
		logger.Errorf("批量插入子请求失败: %v", err)
		return err
	}

	logger.Infof("成功插入 %d 个子请求", len(testObjs))
	return nil
}

// DeleteTestObjectsByRequestID 批量删除子请求
func (s *MongoStore) DeleteTestObjectsByRequestID(ctx context.Context, requestID string) error {
	logger.Infof("开始删除请求 %s 的所有子请求", requestID)
	
	result, err := s.db.Collection(testObjectCollection).DeleteMany(ctx, bson.M{"parent_request_id": requestID})
	if err != nil {
		logger.Errorf("删除子请求失败: %v", err)
		return err
	}

	logger.Infof("成功删除 %d 个子请求", result.DeletedCount)
	return nil
}

// UpdateTestObjectsByRequestID 批量更新子请求
func (s *MongoStore) UpdateTestObjectsByRequestID(ctx context.Context, requestID string, updateData bson.M) error {
	logger.Infof("开始更新请求 %s 的所有子请求", requestID)
	
	result, err := s.db.Collection(testObjectCollection).UpdateMany(
		ctx,
		bson.M{"parent_request_id": requestID},
		bson.M{"$set": updateData},
	)
	if err != nil {
		logger.Errorf("更新子请求失败: %v", err)
		return err
	}

	logger.Infof("成功更新 %d 个子请求", result.ModifiedCount)
	return nil
}

// NewMongoStore 创建新的 MongoDB 存储实例
func NewMongoStore(uri, dbName string) (*MongoStore, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoStore{
		db:     client.Database(dbName),
		client: client,
	}, nil
}

// Close 关闭数据库连接
func (s *MongoStore) Close() error {
	return s.client.Disconnect(context.Background())
} 