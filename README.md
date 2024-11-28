# 🐵 API Monkey Tester

**API Monkey Tester** 是一个用于随机化测试 API 的轻量级工具。用户只需粘贴一个 cURL 格式的接口，便可轻松启动并发和参数变异测试，帮助发现接口潜在的问题。主要解决接口异常参数输入、接口sql注入等潜在性风险问题。为了工具的通用性，
建议接口返回错误使用http错误码，当param错误或者输入异常的时候返回httpcode为非200的参数

---

## 🤬 解决痛点
- **没有质量的输入** ：国内api开发人员的没有责任感和公司没有统一的规范是罪魁祸首
- **渗透攻击**：大部分开发人员并没有足够的安全意识来应对sql、xss等攻击
- **压力测试**：图形化展示并发的测试下接口响应效果


## 🚀 功能特点

- **参数随机化**：自动生成多种参数组合，包括边界值、特殊字符和空值。
- **高并发支持**：模拟多用户请求以测试 API 的稳定性和性能。
- **多协议支持**：兼容 REST 和 GraphQL API。
- **实时响应**：记录每次请求的状态、耗时及错误详情。
- **易于部署**：通过 Docker 或本地运行快速启动。

---

## 📦 安装与运行

### 使用 Docker
```bash
docker-compose up -d
