package adapters

import (
	"context"
	"io"
)

// BackendAdapter defines the unified interface for all simulation backends
// This interface must be implemented by QEMU, Renode, and SkyEye adapters
type BackendAdapter interface {
	// Instance Management
	CreateInstance(ctx context.Context, sessionID string, config *BoardConfig, resources *ResourceConfig) (string, error)
	DestroyInstance(ctx context.Context, instanceID string) error
	
	// Power Management
	PowerOn(ctx context.Context, instanceID string) error
	PowerOff(ctx context.Context, instanceID string) error
	Reset(ctx context.Context, instanceID string) error
	
	// Program Management
	UploadProgram(ctx context.Context, instanceID string, program io.Reader, metadata *ProgramMetadata) (string, error)
	StartProgram(ctx context.Context, instanceID string, programID string, options *StartOptions) error
	PauseProgram(ctx context.Context, instanceID string, programID string) error
	StopProgram(ctx context.Context, instanceID string, programID string) error
	
	// Debug Operations
	SetBreakpoint(ctx context.Context, instanceID string, bp *Breakpoint) error
	RemoveBreakpoint(ctx context.Context, instanceID string, bpID string) error
	StepInstruction(ctx context.Context, instanceID string) error
	Continue(ctx context.Context, instanceID string) error
	
	// State Inspection
	ReadRegisters(ctx context.Context, instanceID string, scope string) (map[string]interface{}, error)
	WriteRegister(ctx context.Context, instanceID string, register string, value interface{}) error
	ReadMemory(ctx context.Context, instanceID string, address uint64, size uint32) ([]byte, error)
	WriteMemory(ctx context.Context, instanceID string, address uint64, data []byte) error
	
	// Snapshot Operations
	CreateSnapshot(ctx context.Context, instanceID string) (string, error)
	RestoreSnapshot(ctx context.Context, instanceID string, snapshotID string) error
	
	// Analysis
	ExportCoverage(ctx context.Context, instanceID string) (string, error)
	ExportTrace(ctx context.Context, instanceID string) (string, error)
	
	// GDB Bridge
	GetGDBServerAddress(ctx context.Context, instanceID string) (string, error)
	
	// Console/Logs
	GetConsoleStream(ctx context.Context, instanceID string) (io.ReadCloser, error)
	
	// Backend Info
	GetBackendType() BackendType
	GetCapabilities() *BackendCapabilities
}

// BackendType represents the type of simulation backend
type BackendType string

const (
	BackendQEMU   BackendType = "qemu"
	BackendRenode BackendType = "renode"
	BackendSkyEye BackendType = "skyeye"
)

// BoardConfig represents the unified board configuration
type BoardConfig struct {
	SystemID     string                 `json:"system_id" yaml:"system_id"`
	Name         string                 `json:"name" yaml:"name"`
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Nodes        []NodeConfig           `json:"nodes" yaml:"nodes"`
	Interconnect *InterconnectConfig    `json:"interconnect,omitempty" yaml:"interconnect,omitempty"`
	Boot         *BootConfig            `json:"boot,omitempty" yaml:"boot,omitempty"`
	Resources    *ResourceConfig        `json:"resources,omitempty" yaml:"resources,omitempty"`
}

// NodeConfig represents a single compute node in the system
type NodeConfig struct {
	ID          string                 `json:"id" yaml:"id"`
	Backend     BackendType            `json:"backend" yaml:"backend"`
	Processor   *ProcessorConfig       `json:"processor" yaml:"processor"`
	Memory      []MemoryRegion         `json:"memory" yaml:"memory"`
	Peripherals []PeripheralConfig     `json:"peripherals" yaml:"peripherals"`
}

// ProcessorConfig represents processor configuration
type ProcessorConfig struct {
	Type        string                 `json:"type" yaml:"type"`        // e.g., "ARM Cortex-M4", "RISC-V RV32"
	Cores       int                    `json:"cores" yaml:"cores"`
	Frequency   uint64                 `json:"frequency" yaml:"frequency"` // Hz
	Features    map[string]interface{} `json:"features,omitempty" yaml:"features,omitempty"`
}

// MemoryRegion represents a memory region
type MemoryRegion struct {
	Type    string `json:"type" yaml:"type"`       // "RAM", "ROM", "Flash"
	Address uint64 `json:"address" yaml:"address"` // Base address
	Size    uint64 `json:"size" yaml:"size"`       // Size in bytes
	Access  string `json:"access" yaml:"access"`   // "RW", "RO", "WO"
}

// PeripheralConfig represents a peripheral device
type PeripheralConfig struct {
	Type       string                 `json:"type" yaml:"type"`
	Name       string                 `json:"name" yaml:"name"`
	Address    uint64                 `json:"address,omitempty" yaml:"address,omitempty"`
	IRQ        []int                  `json:"irq,omitempty" yaml:"irq,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty" yaml:"properties,omitempty"`
}

// InterconnectConfig represents system-level interconnections
type InterconnectConfig struct {
	SharedMemory []SharedMemoryConfig `json:"shared_memory,omitempty" yaml:"shared_memory,omitempty"`
	MMIOMap      []MMIOMapping        `json:"mmio_map,omitempty" yaml:"mmio_map,omitempty"`
	IRQRoutes    []IRQRoute           `json:"irq_routes,omitempty" yaml:"irq_routes,omitempty"`
}

// SharedMemoryConfig represents shared memory between nodes
type SharedMemoryConfig struct {
	ID      string   `json:"id" yaml:"id"`
	Address uint64   `json:"address" yaml:"address"`
	Size    uint64   `json:"size" yaml:"size"`
	Nodes   []string `json:"nodes" yaml:"nodes"` // Node IDs that share this memory
}

// MMIOMapping represents memory-mapped I/O routing
type MMIOMapping struct {
	SourceNode string `json:"source_node" yaml:"source_node"`
	TargetNode string `json:"target_node" yaml:"target_node"`
	Address    uint64 `json:"address" yaml:"address"`
	Size       uint64 `json:"size" yaml:"size"`
}

// IRQRoute represents interrupt routing between nodes
type IRQRoute struct {
	SourceNode string `json:"source_node" yaml:"source_node"`
	SourceIRQ  int    `json:"source_irq" yaml:"source_irq"`
	TargetNode string `json:"target_node" yaml:"target_node"`
	TargetIRQ  int    `json:"target_irq" yaml:"target_irq"`
	Latency    int    `json:"latency,omitempty" yaml:"latency,omitempty"` // cycles
}

// BootConfig represents boot configuration
type BootConfig struct {
	BootROM  string            `json:"bootrom,omitempty" yaml:"bootrom,omitempty"`
	BootArgs []string          `json:"bootargs,omitempty" yaml:"bootargs,omitempty"`
	Env      map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
}

// ResourceConfig represents resource limits
type ResourceConfig struct {
	CPUCores   int    `json:"cpu_cores,omitempty" yaml:"cpu_cores,omitempty"`
	MemoryMB   int    `json:"memory_mb,omitempty" yaml:"memory_mb,omitempty"`
	DiskGB     int    `json:"disk_gb,omitempty" yaml:"disk_gb,omitempty"`
	TimeoutSec int    `json:"timeout_sec,omitempty" yaml:"timeout_sec,omitempty"`
}

// ProgramMetadata represents metadata for uploaded programs
type ProgramMetadata struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"` // "ELF", "BIN", "HEX"
	EntryPoint uint64            `json:"entry_point,omitempty"`
	LoadAddr   uint64            `json:"load_addr,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// StartOptions represents program start options
type StartOptions struct {
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	WaitForGDB  bool              `json:"wait_for_gdb,omitempty"`
	EnableTrace bool              `json:"enable_trace,omitempty"`
}

// Breakpoint represents a debug breakpoint
type Breakpoint struct {
	ID        string `json:"id,omitempty"`
	Address   uint64 `json:"address"`
	Type      string `json:"type"`      // "hardware", "software"
	Condition string `json:"condition,omitempty"`
	Enabled   bool   `json:"enabled"`
}

// BackendCapabilities represents what a backend supports
type BackendCapabilities struct {
	Processors  []string          `json:"processors"`
	Peripherals []string          `json:"peripherals"`
	Buses       []string          `json:"buses"`
	Features    map[string]bool   `json:"features"`
	Limits      map[string]int    `json:"limits"`
}
