package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"med/internal/schemas"
	"net/http"
)

// CreateDoctor   godoc
// @Router       /api/doctor [POST]
// @Security     ApiKeyAuth
// @Summary      Doctor Create
// @Description  Doctor Create
// @Tags         Doctor
// @Produce      json
// @Param        doctor body schemas.DoctorPayload true "Doctor Request"
// @Success      200  {object}  schemas.ResponseSuccess
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) CreateDoctor(c *gin.Context) {
	var params schemas.DoctorPayload
	err := c.ShouldBindJSON(&params)
	if h.handleError(c, err) {
		return
	}

	resp, err := h.services.CreateDoctor(c.Request.Context(), &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateDoctor   godoc
// @Router       /api/doctor/{id} [PUT]
// @Security     ApiKeyAuth
// @Summary      Doctor Update
// @Description  Doctor Update
// @Tags         Doctor
// @Produce      json
// @Param        id path schemas.IDRequest true "Doctor ID"
// @Param        doctor body schemas.DoctorPayload true "Doctor Request"
// @Success      200  {object}  schemas.ResponseSuccess
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) UpdateDoctor(c *gin.Context) {
	var params schemas.DoctorPayload
	err := c.ShouldBindJSON(&params)
	if h.handleError(c, err) {
		return
	}
	id := c.Param("id")
	params.ID = id
	fmt.Println("----------", params.ID)
	fmt.Println("----------", id)

	resp, err := h.services.UpdateDoctor(c.Request.Context(), &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDoctor godoc
// @Security     ApiKeyAuth
// @Router       /api/doctor/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get doctor
// @Description  get doctor
// @Tags         Doctor
// @Accept       json
// @Produce      json
// @Param        id path schemas.IDRequest true "Doctor ID"
// @Success      200  {object}  schemas.Doctor
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) GetDoctor(c *gin.Context) {
	params := schemas.IDRequest{
		ID: c.Param("id"),
	}

	resp, err := h.services.GetDoctor(c.Request.Context(), &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteDoctor   godoc
// @Router       /api/doctor/{id} [delete]
// @Security BearerAuth
// @Summary      Doctor Delete
// @Description  Doctor Delete
// @Tags         Doctor
// @Accept       json
// @Produce      json
// @Param        id path schemas.IDRequest true "Doctor ID"
// @Success      200  {object}  schemas.ResponseSuccess
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) DeleteDoctor(c *gin.Context) {
	ctxWithTimeout, cancel := h.makeContext()
	defer cancel()

	id := c.Param("id")

	resp, err := h.services.DeleteDoctor(ctxWithTimeout, &schemas.IDRequest{ID: id})
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAllDoctors godoc
// @Security     ApiKeyAuth
// @Router       /api/doctors/list [GET]
// @Summary      Get doctors
// @Description  get doctors
// @Tags         Doctor
// @Accept       json
// @Produce      json
// @Param        filter query schemas.GetListRequest false "Filter Request"
// @Success      200  {object}  schemas.ManyDoctors
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) GetAllDoctors(c *gin.Context) {
	ctxWithTimeout, cancel := h.makeContext()
	defer cancel()
	page := cast.ToInt(c.DefaultQuery("page", defaultPage))
	limit := cast.ToInt(c.DefaultQuery("limit", defaultLimit))
	params := schemas.GetListRequest{
		Search: fmt.Sprintf("%%%s%%", c.DefaultQuery("search", "")),
		Page:   int64((page - 1) * limit),
		Limit:  int64(limit),
	}

	resp, err := h.services.GetAllDoctor(ctxWithTimeout, &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}
