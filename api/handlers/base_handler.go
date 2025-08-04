package handlers

import (
	"github.com/ONEST-Network/scheme-manager-adapter/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// BaseHandler contains common handler dependencies
type BaseHandler struct {
	DB     *pgxpool.Pool
	Config *config.Config
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(db *pgxpool.Pool, cfg *config.Config) *BaseHandler {
	return &BaseHandler{DB: db, Config: cfg}
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   errMsg,
	})
}
