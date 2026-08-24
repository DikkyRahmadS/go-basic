package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-basic/internal/modules/user/models"
	"go-basic/internal/modules/user/service"
	"go-basic/internal/pkg/pagination"
	"go-basic/internal/pkg/response"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(
	service service.UserService,
) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Create godoc
// @Summary Create a new user
// @Description Create a new user with the given details
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "User details"
// @Success 201 {object} models.UserResponse "User created successfully"
// @Router /api/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req models.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON request body", nil)
		return
	}

	result, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "User created successfully", result)
}

// FindAll godoc
// @Summary Get all users
// @Description Retrieve a list of all users with optional pagination and filtering
// @Tags users
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param name query string false "Filter by name"
// @Param email query string false "Filter by email"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {array} models.UserResponse "Users retrieved successfully"
// @Router /api/users [get]
func (h *UserHandler) FindAll(c *gin.Context) {
	var req models.FindAllUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", nil)
		return
	}

	result, total, err := h.service.FindAll(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	var opts []response.Option[[]*models.UserResponse]
	if meta := pagination.CalculateMeta(req.Page, req.Limit, total); meta != nil {
		opts = append(opts, response.WithMeta[[]*models.UserResponse](meta))
	}

	response.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", result, opts...)
}

// Update godoc
// @Summary Update an existing user
// @Description Update user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body models.UpdateUserRequest true "User details"
// @Success 200 {object} models.UserResponse "User updated successfully"
// @Router /api/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON request body", nil)
		return
	}

	result, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessResponse(c, http.StatusOK, "User updated successfully", result)
}

// Delete godoc
// @Summary Delete a user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Router /api/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessResponse[any](c, http.StatusOK, "User deleted successfully", nil)
}
