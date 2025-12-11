package adapters

import (
	"context"
	"testing"
)

func TestQEMUAdapter_CreateInstance(t *testing.T) {
	adapter := NewQEMUAdapter("/tmp/test-qemu")
	
	ctx := context.Background()
	sessionID := "test-session-1"
	
	config := &BoardConfig{
		SystemID: "test-system",
		Name:     "Test System",
		Nodes: []NodeConfig{
			{
				ID:      "node1",
				Backend: BackendQEMU,
				Processor: &ProcessorConfig{
					Type:      "ARM Cortex-M4",
					Cores:     1,
					Frequency: 168000000,
				},
				Memory: []MemoryRegion{
					{
						Type:    "RAM",
						Address: 0x20000000,
						Size:    128 * 1024,
						Access:  "RW",
					},
				},
			},
		},
	}
	
	resources := &ResourceConfig{
		CPUCores: 1,
		MemoryMB: 512,
	}
	
	instanceID, err := adapter.CreateInstance(ctx, sessionID, config, resources)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	
	if instanceID == "" {
		t.Fatal("Instance ID should not be empty")
	}
	
	// Clean up
	adapter.DestroyInstance(ctx, instanceID)
}

func TestQEMUAdapter_GetCapabilities(t *testing.T) {
	adapter := NewQEMUAdapter("/tmp/test-qemu")
	
	caps := adapter.GetCapabilities()
	
	if caps == nil {
		t.Fatal("Capabilities should not be nil")
	}
	
	if len(caps.Processors) == 0 {
		t.Error("Should have at least one processor")
	}
	
	if len(caps.Peripherals) == 0 {
		t.Error("Should have at least one peripheral")
	}
	
	if !caps.Features["gdb_support"] {
		t.Error("QEMU should support GDB")
	}
}

func TestQEMUAdapter_BackendType(t *testing.T) {
	adapter := NewQEMUAdapter("/tmp/test-qemu")
	
	backendType := adapter.GetBackendType()
	if backendType != BackendQEMU {
		t.Errorf("Expected backend type QEMU, got %s", backendType)
	}
}

func TestRenodeAdapter_CreateInstance(t *testing.T) {
	adapter := NewRenodeAdapter("/tmp/test-renode")
	
	ctx := context.Background()
	sessionID := "test-session-2"
	
	config := &BoardConfig{
		SystemID: "test-system",
		Name:     "Test System",
	}
	
	resources := &ResourceConfig{
		CPUCores: 1,
		MemoryMB: 512,
	}
	
	instanceID, err := adapter.CreateInstance(ctx, sessionID, config, resources)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	
	if instanceID == "" {
		t.Fatal("Instance ID should not be empty")
	}
	
	// Clean up
	adapter.DestroyInstance(ctx, instanceID)
}

func TestSkyEyeAdapter_CreateInstance(t *testing.T) {
	adapter := NewSkyEyeAdapter("/tmp/test-skyeye")
	
	ctx := context.Background()
	sessionID := "test-session-3"
	
	config := &BoardConfig{
		SystemID: "test-system",
		Name:     "Test System",
	}
	
	resources := &ResourceConfig{
		CPUCores: 1,
		MemoryMB: 512,
	}
	
	instanceID, err := adapter.CreateInstance(ctx, sessionID, config, resources)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	
	if instanceID == "" {
		t.Fatal("Instance ID should not be empty")
	}
	
	// Clean up
	adapter.DestroyInstance(ctx, instanceID)
}
