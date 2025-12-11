package api

import (
	"net/http"

	"github.com/forfire912/virServer/pkg/adapters"
	"github.com/forfire912/virServer/pkg/session"
	"github.com/gin-gonic/gin"
)

// Handler handles API requests
type Handler struct {
	sessionService *session.Service
	adapters       map[adapters.BackendType]adapters.BackendAdapter
}

// NewHandler creates a new API handler
func NewHandler(sessionService *session.Service) *Handler {
	return &Handler{
		sessionService: sessionService,
		adapters:       make(map[adapters.BackendType]adapters.BackendAdapter),
	}
}

// RegisterAdapter registers a backend adapter
func (h *Handler) RegisterAdapter(backend adapters.BackendType, adapter adapters.BackendAdapter) {
	h.adapters[backend] = adapter
}

// GetCapabilities returns backend capabilities
// @Summary Get backend capabilities
// @Description Get list of supported processors, peripherals, and buses with backend support mapping
// @Tags capabilities
// @Produce json
// @Success 200 {object} CapabilitiesResponse
// @Router /capabilities [get]
func (h *Handler) GetCapabilities(c *gin.Context) {
	capabilities := CapabilitiesResponse{
		Processors:  make(map[string][]string),
		Peripherals: make(map[string][]string),
		Buses:       make(map[string][]string),
		Backends:    make(map[string]BackendInfo),
	}
	
	// Collect capabilities from all adapters
	for backendType, adapter := range h.adapters {
		caps := adapter.GetCapabilities()
		backend := string(backendType)
		
		// Map processors to backends
		for _, proc := range caps.Processors {
			capabilities.Processors[proc] = append(capabilities.Processors[proc], backend)
		}
		
		// Map peripherals to backends
		for _, periph := range caps.Peripherals {
			capabilities.Peripherals[periph] = append(capabilities.Peripherals[periph], backend)
		}
		
		// Map buses to backends
		for _, bus := range caps.Buses {
			capabilities.Buses[bus] = append(capabilities.Buses[bus], backend)
		}
		
		// Backend info
		capabilities.Backends[backend] = BackendInfo{
			Type:     backend,
			Features: caps.Features,
			Limits:   caps.Limits,
		}
	}
	
	c.JSON(http.StatusOK, capabilities)
}

// CreateSession creates a new simulation session
// @Summary Create a new session
// @Description Create a new simulation session with specified configuration
// @Tags sessions
// @Accept json
// @Produce json
// @Param request body session.CreateSessionRequest true "Session creation request"
// @Success 201 {object} models.Session
// @Failure 400 {object} ErrorResponse
// @Router /sessions [post]
func (h *Handler) CreateSession(c *gin.Context) {
	var req session.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	sess, err := h.sessionService.CreateSession(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, sess)
}

// GetSession retrieves a session by ID
// @Summary Get session details
// @Description Get details of a specific session
// @Tags sessions
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} models.Session
// @Failure 404 {object} ErrorResponse
// @Router /sessions/{id} [get]
func (h *Handler) GetSession(c *gin.Context) {
	sessionID := c.Param("id")
	
	sess, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "session not found"})
		return
	}
	
	c.JSON(http.StatusOK, sess)
}

// ListSessions lists all sessions
// @Summary List sessions
// @Description List all simulation sessions
// @Tags sessions
// @Produce json
// @Success 200 {array} models.Session
// @Router /sessions [get]
func (h *Handler) ListSessions(c *gin.Context) {
	userID := getUserID(c)
	
	sessions, err := h.sessionService.ListSessions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, sessions)
}

// DeleteSession deletes a session
// @Summary Delete session
// @Description Delete a simulation session
// @Tags sessions
// @Param id path string true "Session ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Router /sessions/{id} [delete]
func (h *Handler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("id")
	
	if err := h.sessionService.DeleteSession(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.Status(http.StatusNoContent)
}

// PowerControl controls power state
// @Summary Control power state
// @Description Power on, off, or reset a session
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body PowerRequest true "Power control request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /sessions/{id}/power [post]
func (h *Handler) PowerControl(c *gin.Context) {
	sessionID := c.Param("id")
	
	var req PowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	if err := h.sessionService.PowerControl(c.Request.Context(), sessionID, req.Action); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Message: "power action completed"})
}

// UploadProgram uploads a program to a session
// @Summary Upload program
// @Description Upload a program (ELF/BIN/HEX) to a session
// @Tags programs
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Session ID"
// @Param file formData file true "Program file"
// @Param name formData string true "Program name"
// @Param type formData string true "Program type (ELF/BIN/HEX)"
// @Success 200 {object} ProgramResponse
// @Failure 400 {object} ErrorResponse
// @Router /sessions/{id}/programs [post]
func (h *Handler) UploadProgram(c *gin.Context) {
	sessionID := c.Param("id")
	
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "file required"})
		return
	}
	
	// TODO: Implement program upload logic
	c.JSON(http.StatusOK, ProgramResponse{
		ID:   "prog-" + sessionID,
		Name: file.Filename,
	})
}

// StartProgram starts a program
// @Summary Start program
// @Description Start execution of an uploaded program
// @Tags programs
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param pid path string true "Program ID"
// @Param request body StartProgramRequest true "Start options"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /sessions/{id}/programs/{pid}/start [post]
func (h *Handler) StartProgram(c *gin.Context) {
	sessionID := c.Param("id")
	programID := c.Param("pid")
	
	var req StartProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = StartProgramRequest{} // Use defaults
	}
	
	adapter, instanceID, err := h.sessionService.GetAdapter(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	
	options := &adapters.StartOptions{
		Args:        req.Args,
		Env:         req.Env,
		WaitForGDB:  req.WaitForGDB,
		EnableTrace: req.EnableTrace,
	}
	
	if err := adapter.StartProgram(c.Request.Context(), instanceID, programID, options); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Message: "program started"})
}

// SetBreakpoint sets a breakpoint
// @Summary Set breakpoint
// @Description Set a debug breakpoint
// @Tags debug
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body adapters.Breakpoint true "Breakpoint"
// @Success 200 {object} adapters.Breakpoint
// @Failure 400 {object} ErrorResponse
// @Router /sessions/{id}/debug/breakpoints [post]
func (h *Handler) SetBreakpoint(c *gin.Context) {
	sessionID := c.Param("id")
	
	var bp adapters.Breakpoint
	if err := c.ShouldBindJSON(&bp); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	adapter, instanceID, err := h.sessionService.GetAdapter(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	
	if err := adapter.SetBreakpoint(c.Request.Context(), instanceID, &bp); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, bp)
}

// ReadRegisters reads register values
// @Summary Read registers
// @Description Read CPU register values
// @Tags debug
// @Produce json
// @Param id path string true "Session ID"
// @Param scope query string false "Register scope"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /sessions/{id}/debug/registers [get]
func (h *Handler) ReadRegisters(c *gin.Context) {
	sessionID := c.Param("id")
	scope := c.DefaultQuery("scope", "general")
	
	adapter, instanceID, err := h.sessionService.GetAdapter(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	
	regs, err := adapter.ReadRegisters(c.Request.Context(), instanceID, scope)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, regs)
}

// Helper function to get user ID from context
func getUserID(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return userID.(string)
}

// Request/Response types

type CapabilitiesResponse struct {
	Processors  map[string][]string   `json:"processors"`
	Peripherals map[string][]string   `json:"peripherals"`
	Buses       map[string][]string   `json:"buses"`
	Backends    map[string]BackendInfo `json:"backends"`
}

type BackendInfo struct {
	Type     string           `json:"type"`
	Features map[string]bool  `json:"features"`
	Limits   map[string]int   `json:"limits"`
}

type PowerRequest struct {
	Action string `json:"action" binding:"required"` // on, off, reset
}

type StartProgramRequest struct {
	Args        []string          `json:"args"`
	Env         map[string]string `json:"env"`
	WaitForGDB  bool              `json:"wait_for_gdb"`
	EnableTrace bool              `json:"enable_trace"`
}

type ProgramResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
