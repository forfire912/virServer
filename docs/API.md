# VirServer API Documentation

## API 端点概览

所有 API 端点都以 `/api/v1` 为前缀。

### 认证

API 支持两种认证方式：

1. **API Key** (推荐用于自动化)
   ```
   X-API-Key: your-api-key
   ```

2. **OAuth2** (推荐用于用户交互)
   ```
   Authorization: Bearer <token>
   ```

## 核心端点

### 1. 能力发现

#### GET /capabilities
获取所有后端支持的处理器、外设和总线类型。

**响应示例：**
```json
{
  "processors": {
    "ARM Cortex-M4": ["qemu", "renode"]
  },
  "peripherals": {
    "UART": ["qemu", "renode", "skyeye"]
  },
  "backends": {
    "qemu": {
      "features": {"gdb_support": true}
    }
  }
}
```

### 2. 会话管理

#### POST /sessions
创建新的仿真会话。

**请求体：**
```json
{
  "name": "会话名称",
  "backend": "qemu|renode|skyeye",
  "board_config": "JSON 配置",
  "board_template": "模板 ID",
  "resources": {
    "cpu_cores": 1,
    "memory_mb": 512
  }
}
```

#### GET /sessions
列出所有会话。

#### GET /sessions/{id}
获取特定会话详情。

#### DELETE /sessions/{id}
删除会话。

### 3. 电源控制

#### POST /sessions/{id}/power
控制会话电源状态。

**请求体：**
```json
{
  "action": "on|off|reset"
}
```

### 4. 程序管理

#### POST /sessions/{id}/programs
上传程序文件（ELF/BIN/HEX）。

**表单数据：**
- `file`: 程序文件
- `name`: 程序名称
- `type`: 文件类型

#### POST /sessions/{id}/programs/{pid}/start
启动程序。

#### POST /sessions/{id}/programs/{pid}/pause
暂停程序。

#### POST /sessions/{id}/programs/{pid}/stop
停止程序。

### 5. 调试

#### POST /sessions/{id}/debug/breakpoints
设置断点。

**请求体：**
```json
{
  "address": 134217728,
  "type": "hardware|software",
  "enabled": true
}
```

#### GET /sessions/{id}/debug/registers
读取寄存器值。

**查询参数：**
- `scope`: general|special|float

#### POST /sessions/{id}/debug/registers/{reg}
写入寄存器值。

#### GET /sessions/{id}/debug/memory
读取内存。

**查询参数：**
- `address`: 内存地址
- `size`: 读取大小

#### POST /sessions/{id}/debug/memory
写入内存。

### 6. 快照

#### POST /sessions/{id}/snapshot
创建快照。

#### POST /sessions/{id}/snapshot/{sid}/restore
恢复快照。

#### GET /sessions/{id}/snapshots
列出所有快照。

### 7. 实时流

#### WebSocket /sessions/{id}/stream
实时控制台和日志流。

**WebSocket 消息格式：**
```json
{
  "type": "stdout|stderr|log",
  "data": "输出内容",
  "timestamp": "ISO8601"
}
```

### 8. 作业管理

#### POST /jobs
创建异步作业。

**请求体：**
```json
{
  "session_id": "会话 ID",
  "type": "coverage|trace|test",
  "options": {}
}
```

#### GET /jobs/{id}
查询作业状态。

#### GET /jobs
列出所有作业。

#### DELETE /jobs/{id}
取消作业。

### 9. 板卡模板

#### GET /templates
列出所有板卡模板。

#### GET /templates/{id}
获取特定模板。

#### POST /templates
创建新模板。

#### PUT /templates/{id}
更新模板。

#### DELETE /templates/{id}
删除模板。

### 10. 模型数据库

#### GET /models/processors
列出所有处理器型号。

#### GET /models/peripherals
列出所有外设类型。

#### GET /models/buses
列出所有总线类型。

## 错误响应

所有错误响应都使用统一格式：

```json
{
  "error": "错误消息描述"
}
```

## 状态码

- `200 OK`: 成功
- `201 Created`: 创建成功
- `204 No Content`: 删除成功
- `400 Bad Request`: 请求错误
- `401 Unauthorized`: 未认证
- `403 Forbidden`: 无权限
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器错误

## 完整 API 规范

访问 Swagger UI 查看完整的交互式 API 文档：
http://localhost:8080/swagger/index.html
