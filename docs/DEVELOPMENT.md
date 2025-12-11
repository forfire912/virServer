# VirServer 开发指南

## 环境设置

### 1. 安装依赖

**Go 环境：**
```bash
# 安装 Go 1.21+
# macOS
brew install go

# Ubuntu/Debian
sudo apt-get install golang-1.21

# 验证安装
go version
```

**其他工具：**
```bash
# Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Docker Compose
sudo apt-get install docker-compose

# Make
sudo apt-get install build-essential
```

### 2. 克隆和构建

```bash
# 克隆仓库
git clone https://github.com/forfire912/virServer.git
cd virServer

# 下载依赖
make deps

# 构建
make build

# 运行测试
make test
```

## 项目结构

```
virServer/
├── cmd/
│   └── server/              # 主程序入口
│       └── main.go
├── pkg/                     # 公共包
│   ├── adapters/            # 后端适配器
│   │   ├── interface.go     # 统一接口定义
│   │   ├── qemu.go          # QEMU 实现
│   │   ├── renode.go        # Renode 实现
│   │   ├── skyeye.go        # SkyEye 实现
│   │   └── adapters_test.go # 测试
│   ├── api/                 # HTTP API
│   │   ├── handlers.go      # 处理器
│   │   └── routes.go        # 路由
│   ├── auth/                # 认证授权
│   ├── board/               # 板卡配置
│   ├── debug/               # 调试服务
│   ├── jobs/                # 作业队列
│   ├── models/              # 数据模型
│   │   └── models.go
│   ├── orchestration/       # 编排服务
│   ├── session/             # 会话管理
│   │   └── service.go
│   ├── sync/                # 同步服务
│   └── analysis/            # 分析服务
├── internal/                # 内部包
│   └── config/              # 配置
│       └── config.go
├── docs/                    # 文档
│   ├── API.md
│   └── ARCHITECTURE.md
├── examples/                # 示例
│   ├── configs/             # 配置示例
│   └── programs/            # 程序示例
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## 开发工作流

### 1. 创建新功能分支

```bash
git checkout -b feature/your-feature-name
```

### 2. 编写代码

遵循 Go 编码规范：
- 使用 `gofmt` 格式化代码
- 添加注释（特别是导出的函数/类型）
- 编写单元测试

### 3. 运行测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./pkg/adapters -v

# 查看覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 4. 本地运行

```bash
# 使用 Docker Compose
docker-compose up -d

# 或直接运行（需要手动启动数据库）
make run

# 查看日志
docker-compose logs -f virserver
```

### 5. 访问 API

```bash
# 健康检查
curl http://localhost:8080/health

# 查询能力
curl http://localhost:8080/api/v1/capabilities

# Swagger UI
open http://localhost:8080/swagger/index.html
```

## 添加新的后端适配器

### 步骤 1: 创建适配器文件

```bash
touch pkg/adapters/mybackend.go
```

### 步骤 2: 实现接口

```go
package adapters

import (
    "context"
    "io"
)

type MyBackendAdapter struct {
    workDir string
}

func NewMyBackendAdapter(workDir string) *MyBackendAdapter {
    return &MyBackendAdapter{
        workDir: workDir,
    }
}

// 实现 BackendAdapter 接口的所有方法
func (a *MyBackendAdapter) CreateInstance(ctx context.Context, sessionID string, config *BoardConfig, resources *ResourceConfig) (string, error) {
    // 实现逻辑
    return "", nil
}

// ... 实现其他方法

func (a *MyBackendAdapter) GetBackendType() BackendType {
    return BackendType("mybackend")
}

func (a *MyBackendAdapter) GetCapabilities() *BackendCapabilities {
    return &BackendCapabilities{
        Processors: []string{"MyProcessor"},
        // ...
    }
}
```

### 步骤 3: 注册适配器

在 `cmd/server/main.go` 中：

```go
// 初始化适配器
mybackendAdapter := adapters.NewMyBackendAdapter(filepath.Join(cfg.Storage.WorkDir, "mybackend"))

// 注册到服务
sessionService.RegisterAdapter(adapters.BackendType("mybackend"), mybackendAdapter)
apiHandler.RegisterAdapter(adapters.BackendType("mybackend"), mybackendAdapter)
```

### 步骤 4: 编写测试

```go
func TestMyBackendAdapter_CreateInstance(t *testing.T) {
    adapter := NewMyBackendAdapter("/tmp/test")
    
    ctx := context.Background()
    config := &BoardConfig{...}
    resources := &ResourceConfig{...}
    
    instanceID, err := adapter.CreateInstance(ctx, "test-session", config, resources)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, instanceID)
}
```

## 添加新的 API 端点

### 步骤 1: 添加处理器

在 `pkg/api/handlers.go` 中：

```go
func (h *Handler) MyNewEndpoint(c *gin.Context) {
    // 解析请求
    var req MyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
        return
    }
    
    // 业务逻辑
    result, err := h.doSomething(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
        return
    }
    
    // 返回响应
    c.JSON(http.StatusOK, result)
}
```

### 步骤 2: 注册路由

在 `pkg/api/routes.go` 中：

```go
v1 := r.Group("/api/v1")
{
    v1.GET("/my-endpoint", handler.MyNewEndpoint)
}
```

### 步骤 3: 添加 Swagger 注释

```go
// MyNewEndpoint does something
// @Summary My endpoint summary
// @Description My endpoint description
// @Tags my-tag
// @Accept json
// @Produce json
// @Param request body MyRequest true "Request body"
// @Success 200 {object} MyResponse
// @Failure 400 {object} ErrorResponse
// @Router /my-endpoint [get]
func (h *Handler) MyNewEndpoint(c *gin.Context) {
    // ...
}
```

### 步骤 4: 重新生成文档

```bash
make swagger
```

## 数据库迁移

### 添加新模型

在 `pkg/models/models.go` 中：

```go
type MyNewModel struct {
    ID        string    `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 注册迁移

在 `cmd/server/main.go` 中：

```go
if err := db.AutoMigrate(
    &models.Session{},
    &models.MyNewModel{},  // 添加新模型
    // ...
); err != nil {
    log.Fatalf("Failed to migrate database: %v", err)
}
```

## 测试指南

### 单元测试

```go
func TestMyFunction(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"
    
    // Act
    result := MyFunction(input)
    
    // Assert
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

### 表驱动测试

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case1", "input1", "output1"},
        {"case2", "input2", "output2"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := MyFunction(tt.input)
            if result != tt.expected {
                t.Errorf("Expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

### 集成测试

```go
func TestAPIIntegration(t *testing.T) {
    // 启动测试服务器
    router := api.SetupRouter(handler)
    
    // 发送请求
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/capabilities", nil)
    router.ServeHTTP(w, req)
    
    // 验证响应
    assert.Equal(t, 200, w.Code)
}
```

## 调试技巧

### 1. 使用 Delve 调试器

```bash
# 安装 delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
dlv debug cmd/server/main.go

# 设置断点
(dlv) break main.main
(dlv) continue
```

### 2. 日志调试

```go
import "log"

log.Printf("Debug: variable = %+v", variable)
```

### 3. 性能分析

```bash
# CPU profiling
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

## 代码风格

### Go 编码规范

1. **命名规范**
   - 包名：小写，简短
   - 导出函数：大写开头
   - 私有函数：小写开头
   - 常量：驼峰命名

2. **注释**
   ```go
   // MyFunction does something useful.
   // It takes a string and returns an error.
   func MyFunction(input string) error {
       // ...
   }
   ```

3. **错误处理**
   ```go
   if err != nil {
       return fmt.Errorf("failed to do something: %w", err)
   }
   ```

4. **接口**
   ```go
   type MyInterface interface {
       DoSomething(ctx context.Context) error
   }
   ```

### 提交规范

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type:**
- feat: 新功能
- fix: 修复
- docs: 文档
- style: 格式
- refactor: 重构
- test: 测试
- chore: 构建/工具

**示例:**
```
feat(adapters): add support for new backend

- Implement BackendAdapter interface
- Add unit tests
- Update documentation

Closes #123
```

## 常见问题

### Q: 如何调试后端适配器？

A: 使用日志和单元测试：
```go
log.Printf("Creating instance: sessionID=%s", sessionID)
```

### Q: 如何处理长时间运行的任务？

A: 使用作业队列系统（正在开发中）。

### Q: 如何添加新的处理器类型？

A: 在数据库的 `processors` 表中添加记录，或在种子数据中添加。

## 资源链接

- **Go 官方文档**: https://golang.org/doc/
- **Gin 文档**: https://gin-gonic.com/docs/
- **GORM 文档**: https://gorm.io/docs/
- **Docker 文档**: https://docs.docker.com/
- **项目 Issues**: https://github.com/forfire912/virServer/issues

## 贡献流程

1. Fork 仓库
2. 创建功能分支
3. 提交代码
4. 编写测试
5. 提交 Pull Request
6. 代码审查
7. 合并

欢迎贡献！🎉
