package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseHandler struct {
	validate *validator.Validate
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		validate: validator.New(),
	}
}

// =====================
// Common helpers
// =====================

// BindAndValidate binds JSON and validates struct
func (b *BaseHandler) BindAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		b.RespondError(c, http.StatusBadRequest, "Invalid request body", err)
		return false
	}
	if err := b.validate.Struct(req); err != nil {
		b.RespondError(c, http.StatusBadRequest, "Validation failed", err)
		return false
	}
	return true
}

// RespondError sends a consistent JSON error
func (b *BaseHandler) RespondError(c *gin.Context, code int, message string, err error) {
	c.JSON(code, gin.H{
		"error":   message,
		"details": err.Error(),
	})
}

// RespondSuccess sends a consistent JSON response
func (b *BaseHandler) RespondSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{"data": data})
}

// GetPaginationParams extracts pagination safely
func (b *BaseHandler) GetPaginationParams(c *gin.Context) (limit, page, offset int) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err = strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	offset = (page - 1) * limit
	return
}
