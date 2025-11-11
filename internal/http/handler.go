package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangira25/user_service/internal/domain"
	"github.com/rangira25/user_service/internal/service"
)

type Handler struct {
	BaseHandler
	svc service.UserService
}

func NewHandler(svc service.UserService) *Handler {
	return &Handler{
		BaseHandler: *NewBaseHandler(),
		svc:         svc,
	}
}

// ========================= HANDLERS =========================

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user account (public or admin creation)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CreateUserReq  true  "User info"
// @Success      201   {object}  domain.UserResp
// @Failure      400   {object}  domain.ErrorResponse
// @Failure      500   {object}  domain.ErrorResponse
// @Router       /users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req domain.CreateUserReq
	if !h.BindAndValidate(c, &req) {
		return
	}

	user, err := h.svc.CreateUser(c.Request.Context(), req)
	if err != nil {
		h.RespondError(c, http.StatusBadRequest, "Failed to create user", err)
		return
	}
	h.RespondSuccess(c, http.StatusCreated, user)
}

// ListUsers godoc
// @Summary      List all users
// @Description  Retrieve a paginated list of users (admin only)
// @Tags         users
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        limit   query     int     false  "Results per page"  example(10)
// @Param        page    query     int     false  "Page number"       example(1)
// @Param        role    query     string  false  "Filter by role"    example("admin")
// @Param        status  query     string  false  "Filter by status"  example("active")
// @Param        q       query     string  false  "Search query"      example("john")
// @Success      200     {object}  domain.ListUsersResponse
// @Failure      500     {object}  domain.ErrorResponse
// @Router       /users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	limit, page, offset := h.GetPaginationParams(c)

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
		h.RespondError(c, http.StatusInternalServerError, "Failed to fetch users", err)
		return
	}

	h.RespondSuccess(c, http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Retrieve a single user by its ID (admin only)
// @Tags         users
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  domain.UserResp
// @Failure      404  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	u, err := h.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		h.RespondError(c, http.StatusNotFound, "User not found", err)
		return
	}
	h.RespondSuccess(c, http.StatusOK, u)
}

// UpdateUser godoc
// @Summary      Update user info
// @Description  Update an existing user's info (admin or self)
// @Tags         users
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "User ID"
// @Param        body  body      domain.UpdateUserReq true  "Updated user info"
// @Success      200   {object}  domain.UserResp
// @Failure      400   {object}  domain.ErrorResponse
// @Failure      500   {object}  domain.ErrorResponse
// @Router       /users/{id} [patch]
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateUserReq
	if !h.BindAndValidate(c, &req) {
		return
	}

	u, err := h.svc.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		h.RespondError(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}
	h.RespondSuccess(c, http.StatusOK, u)
}

// SetStatus godoc
// @Summary      Set user status
// @Description  Change a user's status (active/suspended/blocked) — admin only
// @Tags         users
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id      path      string              true  "User ID"
// @Param        body    body      domain.SetStatusReq true  "Status info"
// @Success      200     {object}  domain.MessageResponse
// @Failure      400     {object}  domain.ErrorResponse
// @Failure      500     {object}  domain.ErrorResponse
// @Router       /users/{id}/status [put]
func (h *Handler) SetStatus(c *gin.Context) {
	id := c.Param("id")
	var req domain.SetStatusReq
	if !h.BindAndValidate(c, &req) {
		return
	}

	if err := h.svc.SetStatus(c.Request.Context(), id, req.Status); err != nil {
		h.RespondError(c, http.StatusInternalServerError, "Failed to update status", err)
		return
	}
	h.RespondSuccess(c, http.StatusOK, domain.MessageResponse{Message: "status updated"})
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Soft delete a user by ID (admin only)
// @Tags         users
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  domain.MessageResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteUser(c.Request.Context(), id); err != nil {
		h.RespondError(c, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}
	h.RespondSuccess(c, http.StatusOK, domain.MessageResponse{Message: "deleted"})
}

// RestoreUser godoc
// @Summary      Restore user
// @Description  Restore a previously deleted user (admin only)
// @Tags         users
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  domain.MessageResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /users/{id}/restore [post]
func (h *Handler) RestoreUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.RestoreUser(c.Request.Context(), id); err != nil {
		h.RespondError(c, http.StatusInternalServerError, "Failed to restore user", err)
		return
	}
	h.RespondSuccess(c, http.StatusOK, domain.MessageResponse{Message: "restored"})
}

// AdminResetPassword godoc
// @Summary      Reset user password (Admin only)
// @Description  Admin resets a user's password to a new one
// @Tags         users
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Param        body body      object  true  "New password"  example({"newPassword": "StrongPassword123"})
// @Success      200  {object}  domain.MessageResponse
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /users/{id}/reset-password [post]
func (h *Handler) AdminResetPassword(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		NewPassword string `json:"newPassword" validate:"required,min=8"`
	}
	if !h.BindAndValidate(c, &body) {
		return
	}

	if err := h.svc.AdminResetPassword(c.Request.Context(), id, body.NewPassword); err != nil {
		h.RespondError(c, http.StatusInternalServerError, "Failed to reset password", err)
		return
	}
	h.RespondSuccess(c, http.StatusOK, domain.MessageResponse{Message: "password reset"})
}

// Login godoc
// @Summary      User login
// @Description  Authenticate a user and return a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.LoginReq  true  "Login credentials"
// @Success      200   {object}  domain.LoginResp
// @Failure      401   {object}  domain.ErrorResponse
// @Failure      500   {object}  domain.ErrorResponse
// @Router       /login [post]
func (h *Handler) Login(c *gin.Context) {
	var req domain.LoginReq
	if !h.BindAndValidate(c, &req) {
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		h.RespondError(c, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}
	h.RespondSuccess(c, http.StatusOK, domain.LoginResp{Token: resp.Token})
}
