package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"erp-internal/internal/modules/user/models"
	"erp-internal/internal/modules/user/service"
	"erp-internal/internal/pkg/pagination"
	"erp-internal/internal/pkg/response"
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

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessResponse[any](c, http.StatusOK, "User deleted successfully", nil)
}
