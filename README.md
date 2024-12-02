# 🐵 API Monkey Tester

**API Monkey Tester** 是一个用于随机化测试 API 的轻量级工具。用户只需粘贴一个 cURL 格式的接口，便可轻松启动并发和参数变异测试，帮助发现接口潜在的问题。主要解决接口异常参数输入、接口sql注入等潜在性风险问题。为了工具的通用性，建议接口返回错误使用http错误码，当param错误或者输入异常的时候返回httpcode为非200的参数。

---

## 🤬 解决痛点
- **没有质量的输入** ：公司没有统一的规范是罪魁祸首
- **渗透攻击**：大部分开发人员并没有足够的安全意识来应对sql、xss等攻击
- **压力测试**：图形化展示并发的测试下接口响应效果

## 🚀 功能特点

- **参数随机化**：自动生成多种参数组合，包括边界值、特殊字符和空值
- **高并发支持**：模拟多用户请求以测试 API 的稳定性和性能
- **多协议支持**：兼容 REST 和 GraphQL API
- **实时响应**：记录每次请求的状态、耗时及错误详情
- **易于部署**：通过 Docker 或本地运行快速启动

## 📦 快速开始

### 使用 Docker 部署

1. **克隆项目**
```bash
git clone https://github.com/yourusername/api-monkey-tester.git
cd api-monkey-tester
```

2. **启动服务**
```bash
chmod +x scripts/start.sh
./scripts/start.sh
```

服务启动后可以访问：
- API 服务: http://localhost:8080
- MongoDB: localhost:27017
- Mongo Express 管理界面: http://localhost:8081

### 使用示例

发送测试请求：
```bash
curl -X POST http://localhost:8080/api/v1/test-requests \
  -H "Content-Type: application/json" \
  -d '{
    "curl": "curl -X POST \"http://example.com/api/users\" -H \"Content-Type: application/json\" -d \"{\\\"name\\\":\\\"john\\\",\\\"age\\\":25}\""
  }'
```

### 查看结果

1. 通过 API 响应直接查看测试结果
2. 访问 http://localhost:8081 使用 Mongo Express 查看详细的测试数据
3. 查看日志：
```bash
docker-compose logs -f api
```

## 🛠️ 项目结构

```
.
├── cmd/                    # 主程序入口
├── conf/                   # 配置文件
├── internal/              # 内部包
│   ├── api/              # API 处理
│   ├── input/            # 输入解析
│   ├── store/            # 数据存储
│   ├── task/             # 任务处理
│   └── logger/           # 日志处理
├── scripts/              # 脚本文件
└── docker-compose.yml    # Docker 编排文件
```

## 🔧 配置说明

配置文件位于 `conf/config.toml`：

```toml
[server]
host = "0.0.0.0"
port = 8080

[database]
uri = "mongodb://admin:password123@mongodb:27017"
database = "monkey_test"

[log]
level = "debug"
filename = "logs/app.log"
maxsize = 100
maxbackups = 10
maxage = 30
compress = true
console = true
```

## 🚫 注意事项

1. 默认配置适用于开发环境，生产环境部署需要：
   - 修改数据库密码
   - 移除 Mongo Express
   - 添加适当的安全配置
   - 配置数据备份策略

2. 测试前请确保目标 API 具备足够的容错能力

## 📝 License

MIT License
