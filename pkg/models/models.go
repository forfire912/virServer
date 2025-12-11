package models

import (
	"time"
)

// Session represents a simulation session
type Session struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Backend     string    `json:"backend"`
	Status      string    `json:"status"`
	BoardConfig string    `json:"board_config" gorm:"type:text"`
	InstanceID  string    `json:"instance_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      string    `json:"user_id"`
}

// SessionStatus represents session status
type SessionStatus string

const (
	SessionCreated  SessionStatus = "created"
	SessionRunning  SessionStatus = "running"
	SessionPaused   SessionStatus = "paused"
	SessionStopped  SessionStatus = "stopped"
	SessionError    SessionStatus = "error"
	SessionDestroyed SessionStatus = "destroyed"
)

// Program represents an uploaded program
type Program struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	SessionID  string    `json:"session_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Path       string    `json:"path"`
	Size       int64     `json:"size"`
	EntryPoint uint64    `json:"entry_point,omitempty"`
	LoadAddr   uint64    `json:"load_addr,omitempty"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// Snapshot represents a simulation snapshot
type Snapshot struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	SessionID  string    `json:"session_id"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Size       int64     `json:"size"`
	CreatedAt  time.Time `json:"created_at"`
	Metadata   string    `json:"metadata" gorm:"type:text"`
}

// Job represents an async job
type Job struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	SessionID   string    `json:"session_id"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Progress    int       `json:"progress"`
	Result      string    `json:"result" gorm:"type:text"`
	ArtifactURL string    `json:"artifact_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// JobStatus represents job status
type JobStatus string

const (
	JobPending   JobStatus = "pending"
	JobRunning   JobStatus = "running"
	JobCompleted JobStatus = "completed"
	JobFailed    JobStatus = "failed"
	JobCancelled JobStatus = "cancelled"
)

// JobType represents job type
type JobType string

const (
	JobCoverage JobType = "coverage"
	JobTrace    JobType = "trace"
	JobTest     JobType = "test"
)

// Processor represents a processor model
type Processor struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Architecture string   `json:"architecture"`
	Vendor      string    `json:"vendor"`
	MaxCores    int       `json:"max_cores"`
	Features    string    `json:"features" gorm:"type:text"`
	Backends    string    `json:"backends"` // Comma-separated list
	CreatedAt   time.Time `json:"created_at"`
}

// Peripheral represents a peripheral model
type Peripheral struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Registers   string    `json:"registers" gorm:"type:text"`
	Backends    string    `json:"backends"` // Comma-separated list
	CreatedAt   time.Time `json:"created_at"`
}

// Bus represents a bus model
type Bus struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Width       int       `json:"width"`
	Backends    string    `json:"backends"` // Comma-separated list
	CreatedAt   time.Time `json:"created_at"`
}

// BoardTemplate represents a board configuration template
type BoardTemplate struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Backend     string    `json:"backend"`
	Config      string    `json:"config" gorm:"type:text"`
	Tags        string    `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// User represents a user in the system
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	APIKey    string    `json:"api_key" gorm:"unique"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Details   string    `json:"details" gorm:"type:text"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
