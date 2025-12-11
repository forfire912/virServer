package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// SetupRouter sets up the API routes
func SetupRouter(handler *Handler) *gin.Engine {
	r := gin.Default()
	
	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	
	// API v1
	v1 := r.Group("/api/v1")
	{
		// Capabilities
		v1.GET("/capabilities", handler.GetCapabilities)
		
		// Sessions
		sessions := v1.Group("/sessions")
		{
			sessions.POST("", handler.CreateSession)
			sessions.GET("", handler.ListSessions)
			sessions.GET("/:id", handler.GetSession)
			sessions.DELETE("/:id", handler.DeleteSession)
			sessions.POST("/:id/power", handler.PowerControl)
			
			// Programs
			sessions.POST("/:id/programs", handler.UploadProgram)
			sessions.POST("/:id/programs/:pid/start", handler.StartProgram)
			sessions.POST("/:id/programs/:pid/pause", handler.PauseProgram)
			sessions.POST("/:id/programs/:pid/stop", handler.StopProgram)
			
			// Debug
			debug := sessions.Group("/:id/debug")
			{
				debug.POST("/breakpoints", handler.SetBreakpoint)
				debug.GET("/registers", handler.ReadRegisters)
				debug.POST("/registers/:reg", handler.WriteRegister)
				debug.GET("/memory", handler.ReadMemory)
				debug.POST("/memory", handler.WriteMemory)
				debug.POST("/step", handler.StepInstruction)
				debug.POST("/continue", handler.Continue)
			}
			
			// Snapshot
			sessions.POST("/:id/snapshot", handler.CreateSnapshot)
			sessions.POST("/:id/snapshot/:sid/restore", handler.RestoreSnapshot)
			sessions.GET("/:id/snapshots", handler.ListSnapshots)
			
			// Console/Logs stream (WebSocket)
			sessions.GET("/:id/stream", handler.StreamConsole)
		}
		
		// Jobs
		jobs := v1.Group("/jobs")
		{
			jobs.POST("", handler.CreateJob)
			jobs.GET("/:id", handler.GetJob)
			jobs.GET("", handler.ListJobs)
			jobs.DELETE("/:id", handler.CancelJob)
		}
		
		// Board templates
		templates := v1.Group("/templates")
		{
			templates.GET("", handler.ListTemplates)
			templates.GET("/:id", handler.GetTemplate)
			templates.POST("", handler.CreateTemplate)
			templates.PUT("/:id", handler.UpdateTemplate)
			templates.DELETE("/:id", handler.DeleteTemplate)
		}
		
		// Model database (processors, peripherals, buses)
		models := v1.Group("/models")
		{
			models.GET("/processors", handler.ListProcessors)
			models.GET("/peripherals", handler.ListPeripherals)
			models.GET("/buses", handler.ListBuses)
		}
	}
	
	return r
}

// Stub handlers for incomplete endpoints

func (h *Handler) PauseProgram(c *gin.Context) {
	c.JSON(200, SuccessResponse{Message: "not implemented"})
}

func (h *Handler) StopProgram(c *gin.Context) {
	c.JSON(200, SuccessResponse{Message: "not implemented"})
}

func (h *Handler) WriteRegister(c *gin.Context) {
	c.JSON(200, SuccessResponse{Message: "not implemented"})
}

func (h *Handler) ReadMemory(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func (h *Handler) WriteMemory(c *gin.Context) {
	c.JSON(200, SuccessResponse{Message: "not implemented"})
}

func (h *Handler) StepInstruction(c *gin.Context) {
	c.JSON(200, SuccessResponse{Message: "not implemented"})
}

func (h *Handler) Continue(c *gin.Context) {
	c.JSON(200, SuccessResponse{Message: "not implemented"})
}

func (h *Handler) CreateSnapshot(c *gin.Context) {
	c.JSON(200, gin.H{"id": "snapshot-1"})
}

func (h *Handler) RestoreSnapshot(c *gin.Context) {
	c.JSON(200, SuccessResponse{Message: "snapshot restored"})
}

func (h *Handler) ListSnapshots(c *gin.Context) {
	c.JSON(200, []interface{}{})
}

func (h *Handler) StreamConsole(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	
	// TODO: Implement console streaming
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *Handler) CreateJob(c *gin.Context) {
	c.JSON(200, gin.H{"id": "job-1"})
}

func (h *Handler) GetJob(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id"), "status": "pending"})
}

func (h *Handler) ListJobs(c *gin.Context) {
	c.JSON(200, []interface{}{})
}

func (h *Handler) CancelJob(c *gin.Context) {
	c.JSON(200, SuccessResponse{Message: "job cancelled"})
}

func (h *Handler) ListTemplates(c *gin.Context) {
	c.JSON(200, []interface{}{})
}

func (h *Handler) GetTemplate(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id")})
}

func (h *Handler) CreateTemplate(c *gin.Context) {
	c.JSON(201, gin.H{"id": "template-1"})
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	c.JSON(200, SuccessResponse{Message: "template updated"})
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	c.Status(204)
}

func (h *Handler) ListProcessors(c *gin.Context) {
	c.JSON(200, []interface{}{})
}

func (h *Handler) ListPeripherals(c *gin.Context) {
	c.JSON(200, []interface{}{})
}

func (h *Handler) ListBuses(c *gin.Context) {
	c.JSON(200, []interface{}{})
}
