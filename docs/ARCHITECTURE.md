# VirServer 架构设计文档

## 1. 系统架构概览

### 1.1 架构层次

```
┌─────────────────────────────────────────────────────────────┐
│                        客户端层                               │
│  (REST Client, gRPC Client, WebSocket Client, GDB Client)   │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                     API Gateway 层                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ REST API     │  │ WebSocket    │  │ Auth         │      │
│  │ (Gin)        │  │ Stream       │  │ Middleware   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                      服务层                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Session      │  │ Model        │  │ Board        │      │
│  │ Service      │  │ Service      │  │ Service      │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Debug        │  │ Job          │  │ Orchestration│      │
│  │ Service      │  │ Service      │  │ Service      │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                    适配器层                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ QEMU         │  │ Renode       │  │ SkyEye       │      │
│  │ Adapter      │  │ Adapter      │  │ Adapter      │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                    仿真后端层                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ QEMU         │  │ Renode       │  │ SkyEye       │      │
│  │ Emulator     │  │ Framework    │  │ Simulator    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                    数据存储层                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ PostgreSQL   │  │ Redis        │  │ MinIO        │      │
│  │ (元数据)      │  │ (缓存/队列)   │  │ (对象存储)    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## 2. 核心组件详解

### 2.1 API Gateway

**职责：**
- HTTP/WebSocket 请求路由
- 认证和授权
- 速率限制
- 日志记录
- OpenAPI 文档服务

**技术实现：**
- Gin Web 框架
- JWT/API Key 认证
- Swagger/OpenAPI 3.0

### 2.2 Session Service

**职责：**
- 仿真会话生命周期管理
- 后端适配器选择和调度
- 资源分配和追踪
- 会话状态持久化

**核心流程：**
```
创建会话 → 选择后端 → 创建实例 → 配置资源 → 保存状态
```

### 2.3 Backend Adapters

**统一接口定义：**
```go
type BackendAdapter interface {
    // 实例管理
    CreateInstance(ctx, sessionID, config, resources) (instanceID, error)
    DestroyInstance(ctx, instanceID) error
    
    // 电源管理
    PowerOn(ctx, instanceID) error
    PowerOff(ctx, instanceID) error
    Reset(ctx, instanceID) error
    
    // 程序管理
    UploadProgram(ctx, instanceID, program, metadata) (programID, error)
    StartProgram(ctx, instanceID, programID, options) error
    
    // 调试操作
    SetBreakpoint(ctx, instanceID, breakpoint) error
    ReadRegisters(ctx, instanceID, scope) (registers, error)
    ReadMemory(ctx, instanceID, address, size) (data, error)
    
    // 快照
    CreateSnapshot(ctx, instanceID) (snapshotID, error)
    RestoreSnapshot(ctx, instanceID, snapshotID) error
    
    // 分析
    ExportCoverage(ctx, instanceID) (artifactURL, error)
    ExportTrace(ctx, instanceID) (artifactURL, error)
    
    // 信息
    GetBackendType() BackendType
    GetCapabilities() *BackendCapabilities
}
```

### 2.4 Board Service

**职责：**
- BoardConfig 模板管理
- 配置验证
- 后端特定配置生成

**BoardConfig 结构：**
```json
{
  "system_id": "唯一标识",
  "name": "系统名称",
  "nodes": [
    {
      "id": "节点ID",
      "backend": "qemu|renode|skyeye",
      "processor": {
        "type": "处理器类型",
        "cores": 1,
        "frequency": 168000000
      },
      "memory": [
        {
          "type": "RAM|ROM|Flash",
          "address": "基地址",
          "size": "大小",
          "access": "RW|RO|WO"
        }
      ],
      "peripherals": [...]
    }
  ],
  "interconnect": {
    "shared_memory": [...],
    "irq_routes": [...],
    "mmio_map": [...]
  }
}
```

## 3. 数据流

### 3.1 会话创建流程

```
1. 客户端 POST /api/v1/sessions
   ↓
2. API Gateway 验证认证
   ↓
3. Session Service 解析 BoardConfig
   ↓
4. 选择/验证后端适配器
   ↓
5. Adapter.CreateInstance()
   ↓
6. 后端启动仿真进程
   ↓
7. 保存会话到数据库
   ↓
8. 返回会话 ID
```

### 3.2 程序执行流程

```
1. 上传程序 POST /sessions/{id}/programs
   ↓
2. 保存到工作目录
   ↓
3. 启动程序 POST /sessions/{id}/programs/{pid}/start
   ↓
4. Adapter.StartProgram()
   ↓
5. 后端加载并执行程序
   ↓
6. WebSocket 流式输出控制台
```

### 3.3 调试流程

```
1. 设置断点 POST /sessions/{id}/debug/breakpoints
   ↓
2. Adapter.SetBreakpoint()
   ↓
3. 后端设置断点
   ↓
4. 程序执行到断点
   ↓
5. 读取寄存器/内存状态
   ↓
6. 返回调试信息
```

## 4. 系统级仿真

### 4.1 多节点协同

对于系统级仿真（多个节点），架构如下：

```
┌──────────────────────────────────────────────┐
│         Orchestration Service                │
│  (管理多节点协调)                             │
└──────────────────────────────────────────────┘
                    │
    ┌───────────────┼───────────────┐
    ↓               ↓               ↓
┌─────────┐    ┌─────────┐    ┌─────────┐
│ Node 1  │    │ Node 2  │    │ Node 3  │
│ (QEMU)  │←──→│(Renode) │←──→│(SkyEye) │
└─────────┘    └─────────┘    └─────────┘
     │              │              │
     └──────────────┼──────────────┘
                    ↓
          ┌──────────────────┐
          │   Sync Service   │
          │ (时间/事件同步)   │
          └──────────────────┘
                    ↓
          ┌──────────────────┐
          │  Memory Proxy    │
          │ (共享内存管理)    │
          └──────────────────┘
```

### 4.2 同步机制

**Timeslice Sync (时间片同步)：**
- 所有节点按固定时间片运行
- Sync Service 协调执行顺序
- 适合松耦合系统

**Step Sync (指令同步)：**
- 节点按指令/周期对齐
- 更精确但性能开销大
- 适合紧密耦合系统

**Event Sync (事件同步)：**
- 基于关键事件推进
- 中断/DMA/总线访问触发同步
- 平衡精度和性能

## 5. 安全架构

### 5.1 认证方式

**API Key：**
```
X-API-Key: <key>
```

**OAuth2：**
```
Authorization: Bearer <token>
```

### 5.2 授权模型

```
User → Role → Permissions → Resources
```

**示例角色：**
- `admin`: 全部权限
- `developer`: 创建/管理自己的会话
- `viewer`: 只读权限

### 5.3 隔离机制

- 每个会话在独立容器/进程中运行
- 资源配额限制（CPU、内存、磁盘）
- 网络隔离
- 文件系统隔离

## 6. 扩展性设计

### 6.1 水平扩展

```
                Load Balancer
                      │
        ┌─────────────┼─────────────┐
        ↓             ↓             ↓
    Server 1      Server 2      Server 3
        │             │             │
        └─────────────┼─────────────┘
                      ↓
              Shared Database
```

### 6.2 后端扩展

添加新后端只需：
1. 实现 `BackendAdapter` 接口
2. 注册到 Session Service
3. 更新能力数据库

### 6.3 功能扩展

- 插件系统（计划中）
- 自定义外设模型（计划中）
- 协议转换器（计划中）

## 7. 监控和运维

### 7.1 指标收集

**系统指标：**
- 活跃会话数
- CPU/内存使用率
- 请求延迟
- 错误率

**业务指标：**
- 会话创建/销毁速率
- 程序执行次数
- 覆盖率生成任务
- 快照大小

### 7.2 日志系统

**日志级别：**
- DEBUG: 详细调试信息
- INFO: 一般操作日志
- WARNING: 警告信息
- ERROR: 错误信息

**审计日志：**
- 用户操作记录
- 资源访问记录
- 配置变更记录

## 8. 部署架构

### 8.1 开发环境

```bash
docker-compose up -d
```

**服务组件：**
- VirServer API (端口 8080)
- PostgreSQL (端口 5432)
- Redis (端口 6379)
- MinIO (端口 9000/9001)

### 8.2 生产环境

**Kubernetes 部署：**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: virserver
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: virserver
        image: virserver:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: postgres-service
```

## 9. API 设计原则

### 9.1 RESTful 风格

- 使用标准 HTTP 方法（GET, POST, PUT, DELETE）
- 资源命名采用名词复数
- 使用 HTTP 状态码表示结果
- 支持分页、过滤、排序

### 9.2 版本控制

- API 版本在 URL 中：`/api/v1/`
- 向后兼容性
- 废弃 API 提前通知

### 9.3 错误处理

统一错误响应格式：
```json
{
  "error": "错误消息",
  "code": "ERROR_CODE",
  "details": {}
}
```

## 10. 性能优化

### 10.1 缓存策略

- Redis 缓存会话状态
- 板卡配置模板缓存
- 能力查询结果缓存

### 10.2 异步处理

- 长时间运行任务使用作业队列
- WebSocket 实时推送
- 后台清理过期会话

### 10.3 资源管理

- 会话超时自动清理
- 磁盘空间监控
- 连接池管理

## 总结

VirServer 采用分层、模块化架构，具有以下特点：

✅ **统一接口**: 屏蔽不同仿真后端差异  
✅ **可扩展**: 易于添加新后端和功能  
✅ **高性能**: 异步处理、缓存优化  
✅ **安全**: 多层认证授权、资源隔离  
✅ **易部署**: Docker/K8s 支持  
✅ **易用**: REST API + Swagger 文档  

系统已实现 MVP 核心功能，为后续扩展奠定了坚实基础。
