package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"monkey-test-api/internal/logger"
)

type MySQLStore struct {
	db *sql.DB
}

// NewMySQLStore 创建新的 MySQL 存储实例
func NewMySQLStore(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// 初始化表结构
	if err := initTables(db); err != nil {
		return nil, err
	}

	return &MySQLStore{db: db}, nil
}

// initTables 初始化数据表
func initTables(db *sql.DB) error {
	// 创建父请求表
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS requests (
			request_id VARCHAR(64) PRIMARY KEY,
			method VARCHAR(10) NOT NULL,
			url TEXT NOT NULL,
			headers JSON,
			body JSON,
			params JSON,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("创建 requests 表失败: %v", err)
	}

	// 创建子请求表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS test_objects (
			test_id VARCHAR(64) PRIMARY KEY,
			parent_request_id VARCHAR(64) NOT NULL,
			method VARCHAR(10) NOT NULL,
			url TEXT NOT NULL,
			headers JSON,
			body JSON,
			params JSON,
			reason TEXT,
			status VARCHAR(20),
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (parent_request_id) REFERENCES requests(request_id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("创建 test_objects 表失败: %v", err)
	}

	return nil
}

// InsertRequest 存储父请求
func (s *MySQLStore) InsertRequest(ctx context.Context, req *Request) error {
	logger.Infof("开始存储请求: %s", req.RequestID)

	headers, err := json.Marshal(req.Headers)
	if err != nil {
		return err
	}

	body, err := json.Marshal(req.Body)
	if err != nil {
		return err
	}

	params, err := json.Marshal(req.Params)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO requests (request_id, method, url, headers, body, params, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, req.RequestID, req.Method, req.URL, headers, body, params, req.Timestamp)

	if err != nil {
		logger.Errorf("存储请求失败: %v", err)
		return err
	}

	logger.Infof("请求存储成功: %s", req.RequestID)
	return nil
}

// GetRequestByID 根据 ID 获取父请求
func (s *MySQLStore) GetRequestByID(ctx context.Context, requestID string) (*Request, error) {
	logger.Debugf("开始查询请求: %s", requestID)

	var req Request
	var headers, body, params []byte

	err := s.db.QueryRowContext(ctx, `
		SELECT request_id, method, url, headers, body, params, timestamp
		FROM requests WHERE request_id = ?
	`, requestID).Scan(&req.RequestID, &req.Method, &req.URL, &headers, &body, &params, &req.Timestamp)

	if err == sql.ErrNoRows {
		logger.Infof("未找到请求: %s", requestID)
		return nil, nil
	}
	if err != nil {
		logger.Errorf("查询请求失败: %v", err)
		return nil, err
	}

	// 解析 JSON 字段
	if err := json.Unmarshal(headers, &req.Headers); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &req.Body); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(params, &req.Params); err != nil {
		return nil, err
	}

	logger.Debugf("成功查询到请求: %s", requestID)
	return &req, nil
}

// GetRequestsByTime 按时间倒序获取父请求
func (s *MySQLStore) GetRequestsByTime(ctx context.Context, limit int64) ([]Request, error) {
	logger.Debugf("开始查询请求列表, limit: %d", limit)

	rows, err := s.db.QueryContext(ctx, `
		SELECT request_id, method, url, headers, body, params, timestamp
		FROM requests ORDER BY timestamp DESC LIMIT ?
	`, limit)
	if err != nil {
		logger.Errorf("查询请求列表失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	var requests []Request
	for rows.Next() {
		var req Request
		var headers, body, params []byte

		if err := rows.Scan(&req.RequestID, &req.Method, &req.URL, &headers, &body, &params, &req.Timestamp); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(headers, &req.Headers); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &req.Body); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(params, &req.Params); err != nil {
			return nil, err
		}

		requests = append(requests, req)
	}

	logger.Debugf("成功查询到 %d 条请求记录", len(requests))
	return requests, nil
}

// InsertTestObjects 批量插入子请求
func (s *MySQLStore) InsertTestObjects(ctx context.Context, testObjs []TestObject) error {
	if len(testObjs) == 0 {
		logger.Debug("没有子请求需要插入")
		return nil
	}

	logger.Infof("开始批量插入 %d 个子请求", len(testObjs))

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO test_objects (
			test_id, parent_request_id, method, url, headers, body, params,
			reason, status, timestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, obj := range testObjs {
		headers, err := json.Marshal(obj.Headers)
		if err != nil {
			return err
		}

		body, err := json.Marshal(obj.Body)
		if err != nil {
			return err
		}

		params, err := json.Marshal(obj.Params)
		if err != nil {
			return err
		}

		_, err = stmt.ExecContext(ctx,
			obj.TestID, obj.ParentRequestID, obj.Method, obj.URL,
			headers, body, params, obj.Reason, obj.Status, obj.Timestamp,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	logger.Infof("成功插入 %d 个子请求", len(testObjs))
	return nil
}

// DeleteTestObjectsByRequestID 批量删除子请求
func (s *MySQLStore) DeleteTestObjectsByRequestID(ctx context.Context, requestID string) error {
	logger.Infof("开始删除请求 %s 的所有子请求", requestID)

	result, err := s.db.ExecContext(ctx, `
		DELETE FROM test_objects WHERE parent_request_id = ?
	`, requestID)
	if err != nil {
		logger.Errorf("��除子请求失败: %v", err)
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	logger.Infof("成功删除 %d 个子请求", count)
	return nil
}

// Close 关闭数据库连接
func (s *MySQLStore) Close() error {
	return s.db.Close()
}

// DeleteRequestByID 删除父请求及其关联数据
func (s *MySQLStore) DeleteRequestByID(ctx context.Context, requestID string) error {
	logger.Infof("开始删除请求及其关联数据: %s", requestID)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除子请求
	if err := s.DeleteTestObjectsByRequestID(ctx, requestID); err != nil {
		return err
	}

	// 删除父请求
	result, err := tx.ExecContext(ctx, `
		DELETE FROM requests WHERE request_id = ?
	`, requestID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return nil
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	logger.Infof("成功删除请求及其关联数据: %s", requestID)
	return nil
}
