# Store 层需求

## 功能描述
Store 层负责与数据存储进行交互，主要包括两个核心操作：
1. **存储用户输入的 cURL 请求**，将其解析为 JSON 格式后存入数据库。
2. **基于解析出的 JSON 请求生成的多个测试对象**，这些对象会基于父子关系通过 `requestID` 进行关联。

通过这两个对象的关系，Store 层将管理请求和测试数据的存储和查询操作。

## 主要功能
1. **存储父请求**
2. **基于requestID 查询父请求**
4. **基于requestID 删除父请求**
5. **基于时间顺序倒序查询父请求**  
  

6. **批量插入子请求，根据requestID**
7. **批量删除子请求，根据requestID**
7. **批量更新子请求，根据requestID**



## 数据库设计

### 数据模型
#### 请求数据模型（父对象）
请求数据存储在 MongoDB 中的文档形式，每个请求会有一个唯一的 `requestID`，以下是请求数据的 JSON 模式：

```json
{
  "requestID": "unique-request-id",  // 唯一的请求标识符
  "method": "GET",                   // 请求方法
  "url": "https://api.example.com",   // 请求URL
  "headers": {                       // 请求头部
    "Content-Type": "application/json"
  },
  "body": "{...}",                   // 请求体 (如果有)
  "params": {                         // 请求参数（查询参数、路径参数等）
    "param1": "value1",
    "param2": "value2"
  },
  "timestamp": "2024-11-29T00:00:00Z"  // 请求时间戳
}

```

#### 测试数据模型（子对象）
每个测试对象会基于请求数据生成，并与父级请求数据通过 requestID 进行关联。以下是测试对象数据的 JSON 模式：
```json
{
  "testID": "unique-test-id",        // 唯一的子请求标识符
  "parentRequestID": "unique-request-id",  // 对应的父请求ID
  "method": "GET",                   // 请求方法（继承自父请求）
  "url": "https://api.example.com",   // 请求URL（继承自父请求）
  "headers": {                       // 请求头部（继承自父请求）
    "Content-Type": "application/json"
  },
  "body": "{...}",                   // 请求体 (如果有，继承自父请求)
  "params": {                         // 子请求的具体参数，可能不同于父请求
    "param1": "value1",
    "param2": "value2"
  },
  "reason": "Test to check invalid parameter handling",  // 测试的原因
  "response": {                      // 子请求的响应数据
    "status": 200,                   // HTTP 状态码
    "body": "{...}",                 // 响应体
    "headers": {                     // 响应头部
      "Content-Type": "application/json"
    },
    "timestamp": "2024-11-29T00:01:00Z"  // 响应时间戳
  },
  "timestamp": "2024-11-29T00:01:00Z"  // 子请求时间戳
}
```

## 数据库操作要求

### 父请求操作

1. **存储父请求**
   存储用户输入的 cURL 请求数据，将其解析为 JSON 格式后存入数据库，生成唯一的 `requestID`，供后续生成子请求使用。

   示例代码：
   ```go
   func (s *Store) InsertRequest(req *Request) error {
       _, err := s.db.Collection("requests").InsertOne(context.Background(), req)
       return err
   }
2. **基于 requestID 查询父请求**
 根据 requestID 查询并获取对应的父请求数据，支持返回请求数据的详细信息。

    示例代码：
   ```go
   func (s *Store) GetRequestByID(requestID string) (*Request, error) {
    var req Request
    err := s.db.Collection("requests").FindOne(context.Background(), bson.M{"request_id": requestID}).Decode(&req)
    if err != nil {
        return nil, err
    }
    return &req, nil}

3. **基于 requestID 删除父请求** 根据 requestID 删除父请求数据及其所有关联的子请求数据，保证数据的完整性。

    示例代码：
   ```go
   func (s *Store) DeleteRequestByID(requestID string) error {
    _, err := s.db.Collection("requests").DeleteOne(context.Background(), bson.M{"request_id": requestID})
    return err}

4. **基于时间顺序倒序查询父请求**
    支持根据时间顺序倒序查询所有父请求数据，方便查看最近的请求。

   示例代码：
    ```go
   func (s *Store) GetRequestsByTime() ([]Request, error) {
    cursor, err := s.db.Collection("requests").Find(context.Background(), bson.M{}, options.Find().SetSort(bson.M{"timestamp": -1}))
    if err != nil {
        return nil, err
    }

    var requests []Request
    if err := cursor.All(context.Background(), &requests); err != nil {
        return nil, err
    }

    return requests, nil}
   
### 子请求操作
1. **批量插入子请求（基于 requestID）**  根据父请求的 requestID 批量插入子请求数据，提高数据插入效率。
   
    示例代码：
    ```go
   func (s *Store) InsertTestObjects(testObjs []TestObject) error {
    var docs []interface{}
    for _, testObj := range testObjs {
        docs = append(docs, testObj)
    }
    _, err := s.db.Collection("test_objects").InsertMany(context.Background(), docs)
    return err}
2. **批量删除子请求（基于 requestID）**  根据父请求的 requestID 批量删除所有关联的子请求，确保数据一致性。

    示例代码：
   ```go
   func (s *Store) DeleteTestObjectsByRequestID(requestID string) error {
    _, err := s.db.Collection("test_objects").DeleteMany(context.Background(), bson.M{"request_id": requestID})
    return err}
3. **批量更新子请求（基于 requestID）**  根据父请求的 requestID 批量更新所有子请求的数据（如状态、重试次数等）。

    示例代码：
   ```go
   func (s *Store) UpdateTestObjectsByRequestID(requestID string, updateData bson.M) error {
    _, err := s.db.Collection("test_objects").UpdateMany(
        context.Background(),
        bson.M{"request_id": requestID},
        bson.M{"$set": updateData},
    )
    return err}
