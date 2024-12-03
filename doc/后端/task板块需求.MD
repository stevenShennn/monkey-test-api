# Task 板块需求文档

## 主要功能

Task 板块负责从父请求生成多个子请求。每个子请求通过父请求中的参数类型和配置文件 `param_test.toml` 来生成不同的测试值，并且每个测试值对应一个子请求。生成的子请求将会存入关系型数据库中，并与父请求建立关联。Task 板块的主要目标是生成和存储测试用例，而不涉及实际的请求发送和响应记录。

## 功能概述

1. **任务输入**：
   - 接收父请求的 JSON 数据，父请求包含用户输入的 `curl` 请求，已拆解为 JSON 格式。
   - 父请求包含若干个参数，这些参数根据其类型（如 `string`, `number`, `boolean` 等）来生成对应的测试用例。

2. **参数变异与组合**：
   - 对父请求中的每个参数，基于 `param_test.toml` 文件中的配置生成多个测试用例。
   - 配置文件中的测试值会根据参数类型（如 `string`，`number`）进行变异，例如：对于 `number` 类型的参数，可以生成多个不同的数值，`string` 类型的参数可能会生成 SQL 注入或 XSS 攻击字符串等。

3. **生成子请求**：
   - 根据父请求中的参数及其类型，生成不同的子请求。
   - 每个子请求包含父请求的所有参数，但有不同的测试值。每个子请求还包含一个 `reason` 字段，用于描述为什么要进行该测试（例如，测试某个边界值，或者测试 SQL 注入攻击等）。

4. **存储子请求**：
   - 子请求的请求数据会存储到关系型数据库中。每个子请求将与父请求通过 `requestID` 建立关联。存储的内容包括子请求的参数、测试值以及测试原因（`reason` 字段）。
   - 存储后，子请求的状态可以标记为“待处理”，后续可以根据需求进行实际请求的发送。

## 数据结构

### 子请求（Child Request）

子请求会继承父请求的所有参数，但参数的值会根据配置文件 `param_test.toml` 中定义的测试值进行变异。每个子请求都需要包含一个 `reason` 字段来说明为什么要进行这个测试。

```json
{
    "requestID": "parent_request_id",   // 关联父请求的ID
    "method": "GET",                     // 请求方法
    "url": "http://example.com/api",     // 请求的URL
    "headers": {
        "Content-Type": "application/json"
    },
    "params": {
        "age": 10000,                   // 子请求的参数
        "name": "testuser"
    },
    "reason": "测试age参数的最大值",     // 测试理由
    "status": "待处理"                   // 请求状态，表示子请求生成后未发出
}


## 配置文件 `param_test.toml` 示例


```TOML
[number]
values = [
    {value = 1, reason = "测试参数%s最小整数值"},
    {value = 100, reason = "测试参数%s普通整数值"},
    {value = 10000, reason = "测试参数%s较大的整数值"},
    {value = 2147483647, reason = "测试参数%sint32最大值"},
    {value = -1, reason = "测试参数%s负整数值"},
    {value = 0.1, reason = "测试参数%s浮动值"}
]

[string]
values = [
    {value = "a", reason = "测试参数%s普通字符串"},
    {value = "test", reason = "测试参数%s普通字符串"},
```


## 任务处理流程

1. **任务接收**：

   Task 板块接收父请求的 JSON 数据，解析其中的参数类型。

2. **加载测试参数值**：

   根据父请求中的每个参数类型，从配置文件 param_test.toml 中加载对应的测试值。
3. **生成子请求**：

   对于父请求中的每个参数，生成多个子请求。每个子请求将包含一个不同的测试值和测试理由（reason 字段）。
   每个子请求将保存父请求的所有参数，但参数值会根据配置文件的测试值进行变异。

4. **存储子请求**：

   生成的子请求会被存储到关系型数据库中，每个子请求都将与父请求通过 requestID 建立关联。
   每个子请求的状态将设置为“待处理”，表示子请求已生成但尚未发出。

5. **任务完成**：

   一旦子请求生成并存储，Task 板块的任务处理流程完成，后续可以通过其他服务或手动触发实际请求的发送。

## MongoDB 数据模型

#### 父请求（Parent Request）

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

#### 子请求（Child Request）

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

## 总结

Task 板块通过接收父请求的数据，基于父请求中的参数类型生成多个子请求。每个子请求的参数值会从配置文件中加载，并进行变异。生成的子请求会存储到 MongoDB 中，并与父请求通过 requestID 建立关联。通过这种方式，Task 板块帮助用户生成全面的 API 测试用例，为后续的实际请求发送和响应分析提供了基础。
