package handlers

import (
	"github.com/gin-gonic/gin"
	"med/internal/auth"
	"med/internal/schemas"
	"med/pkg/utils"
	"net/http"
)

// RegisterUser godoc
// @Router       /api/auth/register [POST]
// @Summary      Register User
// @Description  Registers a new user in the system
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body schemas.SignUpPayload true "User Registration Request"
// @Success      201  {object}  schemas.ResponseSuccess
// @Failure      400  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) RegisterUser(c *gin.Context) {
	var params schemas.SignUpPayload
	err := c.ShouldBindJSON(&params)
	if h.handleError(c, err) {
		return
	}
	err = utils.CheckEmailAndPassword(params.Email)
	if h.handleError(c, err) {
		return
	}
	if !utils.IsValidPassword(params.PasswordHash) {
		c.JSON(400, schemas.ResponseError{"Password is not strong", 400})
		return
	}
	resp, err := h.services.RegisterUser(c.Request.Context(), &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusCreated, schemas.ResponseSuccess{
		Data:   resp,
		Status: 200,
	})
}

// LoginUser godoc
// @Router       /api/auth/login [POST]
// @Summary      Login User
// @Description  Authenticates a user and returns an access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body schemas.SignInPayload true "User Login Credentials"
// @Success      200  {object}  schemas.TokenResponse
// @Failure      400  {object}  schemas.ResponseError
// @Failure      401  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) LoginUser(c *gin.Context) {
	var params schemas.SignInPayload
	err := c.ShouldBindJSON(&params)
	if h.handleError(c, err) {
		return
	}

	resp, err := h.services.LoginUser(c.Request.Context(), &params)
	if h.handleError(c, err) {
		return
	}
	if !utils.IsValidPassword(params.PasswordHash) {
		c.JSON(400, schemas.ResponseError{"Password is not strong", 400})
		return
	}
	jwt, err := auth.GenerateJWTToken(&schemas.TokenPayload{
		ID:   resp.ID,
		Role: "user",
	}, h.log)
	if h.handleError(c, err) {
		return
	}
	c.JSON(200, jwt)
}

// UpdatePassword godoc
// @Router       /api/auth/password [PUT]
// @Security     ApiKeyAuth
// @Summary      Update Password
// @Description  Updates the password for an authenticated user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        password body schemas.ForgetPassPayload true "Password Update Request"
// @Success      200  {object}  schemas.ResponseSuccess
// @Failure      400  {object}  schemas.ResponseError
// @Failure      401  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) UpdatePassword(c *gin.Context) {
	var params schemas.ForgetPassPayload
	err := c.ShouldBindJSON(&params)
	if h.handleError(c, err) {
		return
	}

	resp, err := h.services.UpdatePassword(c.Request.Context(), &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}
