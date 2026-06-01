package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/yourname/yourapp/internal/delivery/http/request"
	"github.com/yourname/yourapp/internal/delivery/http/response"
	"github.com/yourname/yourapp/internal/domain"
	"github.com/yourname/yourapp/internal/usecase/user"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	createUserUC *user.CreateUserUsecase
	getUserUC    *user.GetUserUsecase
	listUsersUC  *user.ListUsersUsecase
	updateUserUC *user.UpdateUserUsecase
	deleteUserUC *user.DeleteUserUsecase
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(
	createUserUC *user.CreateUserUsecase,
	getUserUC *user.GetUserUsecase,
	listUsersUC *user.ListUsersUsecase,
	updateUserUC *user.UpdateUserUsecase,
	deleteUserUC *user.DeleteUserUsecase,
) *UserHandler {
	return &UserHandler{
		createUserUC: createUserUC,
		getUserUC:    getUserUC,
		listUsersUC:  listUsersUC,
		updateUserUC: updateUserUC,
		deleteUserUC: deleteUserUC,
	}
}

// Create handles POST /users
// @Summary Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body request.CreateUserRequest true "Create user request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	input := user.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     req.Role,
	}

	result, err := h.createUserUC.Execute(input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Created(c, result)
}

// Get handles GET /users/:id
// @Summary Get a user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid ID format")
		return
	}

	input := user.GetUserInput{ID: id}
	result, err := h.getUserUC.Execute(input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, result)
}

// List handles GET /users
// @Summary List users with pagination
// @Tags Users
// @Produce json
// @Param name query string false "Filter by name"
// @Param email query string false "Filter by email"
// @Param role query string false "Filter by role"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.Response
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	var req request.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "Invalid query parameters: "+err.Error())
		return
	}

	input := user.ListUsersInput{
		Name:           req.Name,
		Email:          req.Email,
		Role:           req.Role,
		IsActive:       req.IsActive,
		IncludeDeleted: req.IncludeDeleted,
		Page:           req.Page,
		Limit:          req.Limit,
		SortBy:         req.SortBy,
		SortOrder:      req.SortOrder,
	}

	result, err := h.listUsersUC.Execute(input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.SuccessWithMeta(c, result.Users, &response.Meta{
		Page:       result.Pagination.Page,
		Limit:      result.Pagination.Limit,
		Total:      result.Pagination.Total,
		TotalPages: result.Pagination.TotalPages,
	})
}

// Update handles PUT /users/:id
// @Summary Update a user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body request.UpdateUserRequest true "Update user request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid ID format")
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	input := user.UpdateUserInput{
		ID:       id,
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     req.Role,
		IsActive: req.IsActive,
	}

	result, err := h.updateUserUC.Execute(input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "User updated successfully", result)
}

// Delete handles DELETE /users/:id
// @Summary Delete a user
// @Tags Users
// @Param id path string true "User ID"
// @Success 204
// @Failure 404 {object} response.Response
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid ID format")
		return
	}

	input := user.DeleteUserInput{ID: id}
	if err := h.deleteUserUC.Execute(input); err != nil {
		h.handleError(c, err)
		return
	}

	response.NoContent(c)
}

// handleError maps domain errors to HTTP responses
func (h *UserHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, domain.ErrEmailExists):
		response.Conflict(c, err.Error())
	case errors.Is(err, domain.ErrEmailRequired),
		errors.Is(err, domain.ErrInvalidEmail),
		errors.Is(err, domain.ErrNameRequired),
		errors.Is(err, domain.ErrNameTooShort),
		errors.Is(err, domain.ErrPasswordRequired),
		errors.Is(err, domain.ErrPasswordTooShort):
		response.BadRequest(c, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		response.Unauthorized(c, err.Error())
	case errors.Is(err, domain.ErrForbidden):
		response.Forbidden(c, err.Error())
	default:
		c.Error(err) // Log the error
		response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
}
