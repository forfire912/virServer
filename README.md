# VirServer - 统一仿真微服务平台

VirServer 是一个面向多后端（Renode、SkyEye、QEMU）的统一仿真微服务平台，提供一致的 API 接口用于处理器仿真、调试、监控和覆盖率分析。

## 功能特性

### 核心能力
- **多后端支持**: 统一接口支持 QEMU、Renode 和 SkyEye 仿真器
- **能力发现**: 查询支持的处理器、外设和总线类型
- **会话管理**: 创建、列举、查询和删除仿真会话
- **板卡配置**: 灵活的 BoardConfig JSON/YAML 配置，支持模板库
- **仿真控制**: 上电、下电、重置操作
- **程序管理**: 上传、挂载和控制 ELF/BIN/HEX 程序

### 调试功能
- **断点管理**: 设置、删除硬件/软件断点
- **寄存器操作**: 读写 CPU 寄存器
- **内存操作**: 读写内存区域
- **单步调试**: 单指令执行
- **GDB 桥接**: 标准 GDB 协议支持

### 高级功能
- **实时流**: WebSocket 控制台和日志流
- **快照/恢复**: 会话级状态保存和恢复
- **覆盖率分析**: 代码覆盖率收集和导出
- **作业队列**: 异步作业处理（测试、覆盖率、追踪）
- **系统级仿真**: 多节点协同仿真，支持共享内存和中断路由

### 安全与运维
- **身份认证**: API Key 和 OAuth2 支持
- **审计日志**: 完整的操作审计
- **资源配额**: CPU、内存、磁盘限制
- **监控指标**: Prometheus 集成

## 快速开始

### 环境要求
- Go 1.21+
- Docker 和 Docker Compose
- PostgreSQL 15+ (可选，开发环境可使用内存数据库)
- Redis (用于作业队列)
- MinIO 或 S3 (用于存储)

### 使用 Docker Compose 启动

```bash
# 克隆仓库
git clone https://github.com/forfire912/virServer.git
cd virServer

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f virserver
```

服务将在以下端口启动：
- **VirServer API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **PostgreSQL**: localhost:5432
- **MinIO Console**: http://localhost:9001
- **Redis**: localhost:6379

### 本地开发

```bash
# 安装依赖
make deps

# 运行测试
make test

# 构建
make build

# 运行服务器
make run
```

## API 使用示例

### 1. 查询后端能力

```bash
curl http://localhost:8080/api/v1/capabilities
```

### 2. 创建仿真会话

```bash
curl -X POST http://localhost:8080/api/v1/sessions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My STM32 Session",
    "backend": "qemu",
    "board_template": "stm32f4-disco",
    "resources": {"cpu_cores": 1, "memory_mb": 512}
  }'
```

### 3. 控制电源

```bash
curl -X POST http://localhost:8080/api/v1/sessions/{session_id}/power \
  -H "Content-Type: application/json" \
  -d '{"action": "on"}'
```

更多 API 示例请访问 Swagger UI: http://localhost:8080/swagger/index.html

## 架构概述

```
API Gateway → Session Manager → Backend Adapters (QEMU/Renode/SkyEye)
```

完整架构文档请参见 `docs/architecture.md`

## 项目结构

```
virServer/
├── cmd/server/          # 主程序入口
├── pkg/
│   ├── adapters/        # 后端适配器
│   ├── api/             # HTTP API
│   ├── session/         # 会话管理
│   └── models/          # 数据模型
├── examples/            # 示例配置
└── docker-compose.yml   # Docker 配置
```

## 开发路线图

### Phase 1: MVP (v0.1.0) ✓
- [x] 核心架构和接口
- [x] 会话管理
- [x] 基础 API
- [x] BoardConfig 支持

### Phase 2-4: 详细规划
请参见项目 Issues 和 Milestones

## 许可证

Apache 2.0 License
