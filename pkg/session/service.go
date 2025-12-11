package session

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/forfire912/virServer/pkg/adapters"
	"github.com/forfire912/virServer/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service manages simulation sessions
type Service struct {
	db          *gorm.DB
	mu          sync.RWMutex
	adapters    map[adapters.BackendType]adapters.BackendAdapter
	sessions    map[string]*SessionRuntime
}

// SessionRuntime holds runtime information for a session
type SessionRuntime struct {
	Session    *models.Session
	Adapter    adapters.BackendAdapter
	InstanceID string
}

// NewService creates a new session service
func NewService(db *gorm.DB) *Service {
	return &Service{
		db:       db,
		adapters: make(map[adapters.BackendType]adapters.BackendAdapter),
		sessions: make(map[string]*SessionRuntime),
	}
}

// RegisterAdapter registers a backend adapter
func (s *Service) RegisterAdapter(backend adapters.BackendType, adapter adapters.BackendAdapter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.adapters[backend] = adapter
}

// CreateSession creates a new simulation session
func (s *Service) CreateSession(ctx context.Context, req *CreateSessionRequest) (*models.Session, error) {
	// Validate backend
	backend := adapters.BackendType(req.Backend)
	if backend == "" {
		backend = adapters.BackendQEMU // Default
	}
	
	s.mu.RLock()
	adapter, exists := s.adapters[backend]
	s.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("backend not supported: %s", backend)
	}
	
	// Parse BoardConfig
	var boardConfig adapters.BoardConfig
	if req.BoardConfig != "" {
		if err := json.Unmarshal([]byte(req.BoardConfig), &boardConfig); err != nil {
			return nil, fmt.Errorf("invalid board config: %w", err)
		}
	} else if req.BoardTemplate != "" {
		// Load template
		var template models.BoardTemplate
		if err := s.db.Where("id = ?", req.BoardTemplate).First(&template).Error; err != nil {
			return nil, fmt.Errorf("template not found: %w", err)
		}
		if err := json.Unmarshal([]byte(template.Config), &boardConfig); err != nil {
			return nil, fmt.Errorf("invalid template config: %w", err)
		}
	} else {
		return nil, fmt.Errorf("either board_config or board_template required")
	}
	
	// Create session record
	session := &models.Session{
		ID:      uuid.New().String(),
		Name:    req.Name,
		Backend: string(backend),
		Status:  string(models.SessionCreated),
		UserID:  getUserIDFromContext(ctx),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	configBytes, _ := json.Marshal(boardConfig)
	session.BoardConfig = string(configBytes)
	
	// Create backend instance
	resources := &adapters.ResourceConfig{
		CPUCores:   req.Resources.CPUCores,
		MemoryMB:   req.Resources.MemoryMB,
		DiskGB:     req.Resources.DiskGB,
		TimeoutSec: req.Resources.TimeoutSec,
	}
	
	instanceID, err := adapter.CreateInstance(ctx, session.ID, &boardConfig, resources)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}
	
	session.InstanceID = instanceID
	
	// Save to database
	if err := s.db.Create(session).Error; err != nil {
		// Cleanup instance on error
		adapter.DestroyInstance(ctx, instanceID)
		return nil, fmt.Errorf("failed to save session: %w", err)
	}
	
	// Store runtime info
	s.mu.Lock()
	s.sessions[session.ID] = &SessionRuntime{
		Session:    session,
		Adapter:    adapter,
		InstanceID: instanceID,
	}
	s.mu.Unlock()
	
	return session, nil
}

// GetSession retrieves a session by ID
func (s *Service) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	var session models.Session
	if err := s.db.Where("id = ?", sessionID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// ListSessions lists all sessions for a user
func (s *Service) ListSessions(ctx context.Context, userID string) ([]*models.Session, error) {
	var sessions []*models.Session
	query := s.db.Order("created_at DESC")
	
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	
	if err := query.Find(&sessions).Error; err != nil {
		return nil, err
	}
	
	return sessions, nil
}

// DeleteSession deletes a session
func (s *Service) DeleteSession(ctx context.Context, sessionID string) error {
	s.mu.Lock()
	runtime, exists := s.sessions[sessionID]
	if exists {
		// Destroy backend instance
		runtime.Adapter.DestroyInstance(ctx, runtime.InstanceID)
		delete(s.sessions, sessionID)
	}
	s.mu.Unlock()
	
	// Delete from database
	if err := s.db.Where("id = ?", sessionID).Delete(&models.Session{}).Error; err != nil {
		return err
	}
	
	return nil
}

// PowerControl controls the power state of a session
func (s *Service) PowerControl(ctx context.Context, sessionID string, action string) error {
	s.mu.RLock()
	runtime, exists := s.sessions[sessionID]
	s.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}
	
	var err error
	switch action {
	case "on":
		err = runtime.Adapter.PowerOn(ctx, runtime.InstanceID)
		if err == nil {
			s.updateSessionStatus(sessionID, models.SessionRunning)
		}
	case "off":
		err = runtime.Adapter.PowerOff(ctx, runtime.InstanceID)
		if err == nil {
			s.updateSessionStatus(sessionID, models.SessionStopped)
		}
	case "reset":
		err = runtime.Adapter.Reset(ctx, runtime.InstanceID)
	default:
		return fmt.Errorf("invalid power action: %s", action)
	}
	
	return err
}

// GetAdapter returns the adapter for a session
func (s *Service) GetAdapter(sessionID string) (adapters.BackendAdapter, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	runtime, exists := s.sessions[sessionID]
	if !exists {
		return nil, "", fmt.Errorf("session not found: %s", sessionID)
	}
	
	return runtime.Adapter, runtime.InstanceID, nil
}

// Helper function to update session status
func (s *Service) updateSessionStatus(sessionID string, status models.SessionStatus) {
	s.db.Model(&models.Session{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
		"status":     string(status),
		"updated_at": time.Now(),
	})
}

// Helper function to get user ID from context
func getUserIDFromContext(ctx context.Context) string {
	userID := ctx.Value("user_id")
	if userID != nil {
		return userID.(string)
	}
	return "anonymous"
}

// CreateSessionRequest represents a request to create a session
type CreateSessionRequest struct {
	Name          string         `json:"name" binding:"required"`
	Backend       string         `json:"backend"`
	BoardConfig   string         `json:"board_config"`
	BoardTemplate string         `json:"board_template"`
	Resources     ResourceConfig `json:"resources"`
}

// ResourceConfig represents resource configuration
type ResourceConfig struct {
	CPUCores   int `json:"cpu_cores"`
	MemoryMB   int `json:"memory_mb"`
	DiskGB     int `json:"disk_gb"`
	TimeoutSec int `json:"timeout_sec"`
}
