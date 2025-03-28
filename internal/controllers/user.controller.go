package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"med/internal/schemas"
	"net/http"
)

// UpdateUser godoc
// @Router       /api/user/{id} [PUT]
// @Security     ApiKeyAuth
// @Summary      Update User Profile
// @Description  Updates the profile of an existing user
// @Tags         User
// @Produce      json
// @Param        id path string true "User ID"
// @Param        user body schemas.UpdateUserProfilePayload true "User Update Request"
// @Success      200  {object}  schemas.IDResponse
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) UpdateUser(c *gin.Context) {
	var params schemas.UpdateUserProfilePayload
	err := c.ShouldBindJSON(&params)
	if h.handleError(c, err) {
		return
	}
	params.ID = c.Param("id")

	resp, err := h.services.UpdateUser(c.Request.Context(), &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteUser godoc
// @Router       /api/user/{id} [DELETE]
// @Security     ApiKeyAuth
// @Summary      Delete User
// @Description  Deletes an existing user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200  {object}  schemas.ResponseSuccess
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) DeleteUser(c *gin.Context) {
	ctxWithTimeout, cancel := h.makeContext()
	defer cancel()

	params := schemas.IDRequest{
		ID: c.Param("id"),
	}

	resp, err := h.services.DeleteUser(ctxWithTimeout, &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUser godoc
// @Router       /api/user/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get User
// @Description  Retrieves a user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200  {object}  schemas.UserResponse
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) GetUser(c *gin.Context) {
	params := schemas.IDRequest{
		ID: c.Param("id"),
	}

	resp, err := h.services.GetUser(c.Request.Context(), &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAllUsers godoc
// @Router       /api/users/list [GET]
// @Security     ApiKeyAuth
// @Summary      Get All Users
// @Description  Retrieves a list of all users with optional search filter
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        search query string false "Search term"
// @Success      200  {object}  schemas.ManyUsers
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) GetAllUsers(c *gin.Context) {
	ctxWithTimeout, cancel := h.makeContext()
	defer cancel()

	params := schemas.GetSearchRequest{
		Search: fmt.Sprintf("%%%s%%", c.DefaultQuery("search", "")),
	}

	resp, err := h.services.GetAllUsers(ctxWithTimeout, &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}
