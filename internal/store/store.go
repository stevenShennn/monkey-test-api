package store

import (
	"context"
)

// Store 定义存储层接口
type Store interface {
	// 父请求相关操作
	InsertRequest(ctx context.Context, req *Request) error
	GetRequestByID(ctx context.Context, requestID string) (*Request, error)
	GetRequestsByTime(ctx context.Context, limit int64) ([]Request, error)
	DeleteRequestByID(ctx context.Context, requestID string) error

	// 子请求相关操作
	InsertTestObjects(ctx context.Context, testObjs []TestObject) error
	DeleteTestObjectsByRequestID(ctx context.Context, requestID string) error

	// 关闭连接
	Close() error
} 