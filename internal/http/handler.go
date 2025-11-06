package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rangira25/user_service/internal/domain"
	"github.com/rangira25/user_service/internal/service"
)

type Handler struct {
	svc    service.UserService
	validate *validator.Validate
}

func NewHandler(svc service.UserService) *Handler {
	v := validator.New()
	return &Handler{svc: svc, validate: v}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req domain.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h *Handler) ListUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	offset := (page - 1) * limit

	filter := map[string]interface{}{}
	if s := c.Query("status"); s != "" {
		filter["status"] = s
	}
	if r := c.Query("role"); r != "" {
		filter["role"] = r
	}
	if q := c.Query("q"); q != "" {
		filter["search"] = q
	}

	users, total, err := h.svc.ListUsers(c.Request.Context(), filter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users, "total": total, "page": page, "limit": limit})
}

func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	u, err := h.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": u})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	u, err := h.svc.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"data": u})
}

func (h *Handler) SetStatus(c *gin.Context) {
	id := c.Param("id")
	var req domain.SetStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if err := h.svc.SetStatus(c.Request.Context(), id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) RestoreUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.RestoreUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"message": "restored"})
}

func (h *Handler) AdminResetPassword(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		NewPassword string `json:"newPassword" validate:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if err := h.validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if err := h.svc.AdminResetPassword(c.Request.Context(), id, body.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password reset"})
}
