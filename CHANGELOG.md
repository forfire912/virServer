# Changelog

All notable changes to VirServer will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-12-11

### Added

#### Core Infrastructure
- ✅ Go-based microservice architecture
- ✅ Unified `BackendAdapter` interface for simulation backends
- ✅ BoardConfig schema (JSON/YAML) with validation
- ✅ API Gateway with Gin web framework
- ✅ OpenAPI/Swagger documentation support
- ✅ Configuration management with environment variables

#### Backend Adapters
- ✅ QEMU adapter with instance management
- ✅ Renode adapter skeleton implementation
- ✅ SkyEye adapter skeleton implementation
- ✅ Capability discovery for each backend
- ✅ Power control (on/off/reset) interface
- ✅ GDB server address support

#### Data Models
- ✅ Session model with lifecycle tracking
- ✅ Program model for uploaded binaries
- ✅ Snapshot model for state persistence
- ✅ Job model for async operations
- ✅ Processor/Peripheral/Bus capability models
- ✅ Board template model
- ✅ User and audit log models
- ✅ GORM integration with PostgreSQL/SQLite

#### API Endpoints (30+)
- ✅ `GET /health` - Health check
- ✅ `GET /api/v1/capabilities` - Backend capability discovery
- ✅ `POST /api/v1/sessions` - Create session
- ✅ `GET /api/v1/sessions` - List sessions
- ✅ `GET /api/v1/sessions/{id}` - Get session details
- ✅ `DELETE /api/v1/sessions/{id}` - Delete session
- ✅ `POST /api/v1/sessions/{id}/power` - Power control
- ✅ `POST /api/v1/sessions/{id}/programs` - Upload program
- ✅ `POST /api/v1/sessions/{id}/programs/{pid}/start` - Start program
- ✅ `POST /api/v1/sessions/{id}/programs/{pid}/pause` - Pause program
- ✅ `POST /api/v1/sessions/{id}/programs/{pid}/stop` - Stop program
- ✅ `POST /api/v1/sessions/{id}/debug/breakpoints` - Set breakpoint
- ✅ `GET /api/v1/sessions/{id}/debug/registers` - Read registers
- ✅ `POST /api/v1/sessions/{id}/debug/registers/{reg}` - Write register
- ✅ `GET /api/v1/sessions/{id}/debug/memory` - Read memory
- ✅ `POST /api/v1/sessions/{id}/debug/memory` - Write memory
- ✅ `POST /api/v1/sessions/{id}/debug/step` - Step instruction
- ✅ `POST /api/v1/sessions/{id}/debug/continue` - Continue execution
- ✅ `POST /api/v1/sessions/{id}/snapshot` - Create snapshot
- ✅ `POST /api/v1/sessions/{id}/snapshot/{sid}/restore` - Restore snapshot
- ✅ `GET /api/v1/sessions/{id}/snapshots` - List snapshots
- ✅ `GET /api/v1/sessions/{id}/stream` - WebSocket console stream
- ✅ `POST /api/v1/jobs` - Create async job
- ✅ `GET /api/v1/jobs/{id}` - Get job status
- ✅ `GET /api/v1/jobs` - List jobs
- ✅ `DELETE /api/v1/jobs/{id}` - Cancel job
- ✅ Template management endpoints
- ✅ Model database query endpoints

#### Services
- ✅ Session management service with adapter routing
- ✅ Model service for capability database
- ✅ Board service framework for config management
- ✅ Authentication middleware framework
- ✅ WebSocket support for real-time streaming

#### Deployment & DevOps
- ✅ Dockerfile for containerization
- ✅ Docker Compose with PostgreSQL, Redis, MinIO
- ✅ Makefile with build/test/clean targets
- ✅ GitHub Actions CI pipeline
- ✅ Automated testing on push/PR
- ✅ SQLite fallback for development

#### Documentation
- ✅ Comprehensive README with quickstart
- ✅ API documentation with examples
- ✅ Architecture design document
- ✅ Development guide
- ✅ Board configuration examples (STM32F4, multi-node)
- ✅ Swagger UI integration

#### Testing
- ✅ Unit tests for backend adapters
- ✅ Test coverage tracking
- ✅ CI integration for automated testing

#### Examples
- ✅ STM32F4 Discovery board configuration
- ✅ Multi-node heterogeneous system example
- ✅ API usage examples in documentation

### Technical Details

**Languages & Frameworks:**
- Go 1.21
- Gin Web Framework
- GORM ORM
- Gorilla WebSocket

**Database:**
- PostgreSQL 15 (production)
- SQLite (development fallback)

**Infrastructure:**
- Docker & Docker Compose
- Redis for caching/queues
- MinIO for object storage

**Architecture:**
- Layered architecture (API → Service → Adapter → Backend)
- Unified adapter interface pattern
- RESTful API design
- WebSocket for real-time streaming

### Known Limitations

- 🔄 Debug operations are framework only (not fully implemented)
- 🔄 Program upload/execution needs full implementation
- 🔄 GDB bridging needs implementation
- 🔄 Coverage/trace collection needs implementation
- 🔄 Job queue system needs implementation
- 🔄 Sync service for multi-node coordination not implemented
- 🔄 Memory proxy not implemented
- 🔄 Interrupt bridge not implemented

### Security

- ✅ API Key authentication framework
- ✅ OAuth2 support framework
- ✅ Audit logging models
- ⚠️ RBAC not fully implemented

## [Unreleased]

### Planned for v0.2.0

#### Features
- Full GDB protocol bridging implementation
- Complete program upload and execution logic
- Coverage collection and export (lcov/gcov)
- Trace collection and export
- Job queue with Redis backend
- Prometheus metrics integration
- Multi-node orchestration service

#### Improvements
- Enhanced error handling
- Better logging
- Performance optimization
- Resource quota enforcement
- Connection pooling

#### Documentation
- Integration test examples
- Deployment guides for Kubernetes
- Security best practices
- API versioning strategy

### Planned for v0.3.0

#### Features
- Sync service for multi-node coordination
- Memory proxy for shared memory
- Interrupt bridge for cross-node IRQs
- Hot-pluggable peripheral support
- Virtual network between nodes
- Performance counter support

### Planned for v1.0.0

#### Features
- Production-ready security (RBAC)
- Horizontal scaling support
- High availability configuration
- Advanced monitoring and alerting
- Complete test coverage (>80%)
- Performance benchmarks
- Migration tools

## Version History

- **v0.1.0** (2025-12-11) - Initial MVP release with core architecture

---

## How to Read This Changelog

- ✅ Completed feature
- 🔄 In progress / Partial implementation
- ⚠️ Needs attention / Known issue
- 🎯 Planned feature

## Contributing

See [DEVELOPMENT.md](docs/DEVELOPMENT.md) for contribution guidelines.

## License

Apache 2.0 - See LICENSE file for details
