package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/forfire912/virServer/internal/config"
	"github.com/forfire912/virServer/pkg/adapters"
	"github.com/forfire912/virServer/pkg/api"
	"github.com/forfire912/virServer/pkg/models"
	"github.com/forfire912/virServer/pkg/session"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title VirServer API
// @version 1.0
// @description Unified Simulation Microservice Platform for Multiple Backends (QEMU, Renode, SkyEye)
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://github.com/forfire912/virServer
// @contact.email support@virserver.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

var (
	Version   = "0.1.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)
	
	// Print version info
	log.Printf("VirServer v%s (build: %s, commit: %s)", Version, BuildTime, GitCommit)
	
	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.Session{},
		&models.Program{},
		&models.Snapshot{},
		&models.Job{},
		&models.Processor{},
		&models.Peripheral{},
		&models.Bus{},
		&models.BoardTemplate{},
		&models.User{},
		&models.AuditLog{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	
	// Create work directories
	createDirectories(cfg)
	
	// Initialize services
	sessionService := session.NewService(db)
	
	// Initialize and register backend adapters
	qemuAdapter := adapters.NewQEMUAdapter(filepath.Join(cfg.Storage.WorkDir, "qemu"))
	renodeAdapter := adapters.NewRenodeAdapter(filepath.Join(cfg.Storage.WorkDir, "renode"))
	skyeyeAdapter := adapters.NewSkyEyeAdapter(filepath.Join(cfg.Storage.WorkDir, "skyeye"))
	
	sessionService.RegisterAdapter(adapters.BackendQEMU, qemuAdapter)
	sessionService.RegisterAdapter(adapters.BackendRenode, renodeAdapter)
	sessionService.RegisterAdapter(adapters.BackendSkyEye, skyeyeAdapter)
	
	// Initialize API handler
	apiHandler := api.NewHandler(sessionService)
	apiHandler.RegisterAdapter(adapters.BackendQEMU, qemuAdapter)
	apiHandler.RegisterAdapter(adapters.BackendRenode, renodeAdapter)
	apiHandler.RegisterAdapter(adapters.BackendSkyEye, skyeyeAdapter)
	
	// Setup routes
	router := api.SetupRouter(apiHandler)
	
	// Seed initial data
	seedInitialData(db)
	
	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	log.Printf("Swagger UI available at http://%s/swagger/index.html", addr)
	
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// If postgres connection fails, use in-memory SQLite for development
		log.Printf("Warning: PostgreSQL connection failed, using in-memory SQLite: %v", err)
		
		// Import SQLite driver
		sqliteDB, sqliteErr := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if sqliteErr != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", sqliteErr)
		}
		log.Println("Using in-memory SQLite database for development")
		return sqliteDB, nil
	}
	
	return db, nil
}

func createDirectories(cfg *config.Config) {
	dirs := []string{
		cfg.Storage.ArtifactPath,
		cfg.Storage.SnapshotPath,
		cfg.Storage.WorkDir,
		filepath.Join(cfg.Storage.WorkDir, "qemu"),
		filepath.Join(cfg.Storage.WorkDir, "renode"),
		filepath.Join(cfg.Storage.WorkDir, "skyeye"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warning: Failed to create directory %s: %v", dir, err)
		}
	}
}

func seedInitialData(db *gorm.DB) {
	// Seed processors
	processors := []models.Processor{
		{
			ID:           "cortex-m3",
			Name:         "ARM Cortex-M3",
			Type:         "ARM",
			Architecture: "ARMv7-M",
			Vendor:       "ARM",
			MaxCores:     1,
			Backends:     "qemu,renode,skyeye",
		},
		{
			ID:           "cortex-m4",
			Name:         "ARM Cortex-M4",
			Type:         "ARM",
			Architecture: "ARMv7E-M",
			Vendor:       "ARM",
			MaxCores:     1,
			Backends:     "qemu,renode",
		},
		{
			ID:           "rv32",
			Name:         "RISC-V RV32",
			Type:         "RISC-V",
			Architecture: "RV32IMAC",
			Vendor:       "RISC-V",
			MaxCores:     4,
			Backends:     "qemu,renode",
		},
	}
	
	for _, proc := range processors {
		db.FirstOrCreate(&proc, models.Processor{ID: proc.ID})
	}
	
	// Seed board templates
	templates := []models.BoardTemplate{
		{
			ID:          "stm32f4-disco",
			Name:        "STM32F4 Discovery",
			Description: "STM32F4 Discovery board with Cortex-M4",
			Backend:     "qemu",
			Config:      `{"system_id":"stm32f4-disco","name":"STM32F4 Discovery","nodes":[{"id":"mcu","backend":"qemu","processor":{"type":"ARM Cortex-M4","cores":1,"frequency":168000000},"memory":[{"type":"Flash","address":134217728,"size":1048576,"access":"RX"},{"type":"RAM","address":536870912,"size":196608,"access":"RW"}],"peripherals":[{"type":"UART","name":"USART2","address":1073759232},{"type":"GPIO","name":"GPIOD","address":1073889280}]}]}`,
			Tags:        "stm32,arm,cortex-m4",
		},
	}
	
	for _, tmpl := range templates {
		db.FirstOrCreate(&tmpl, models.BoardTemplate{ID: tmpl.ID})
	}
	
	log.Println("Initial data seeded successfully")
}
