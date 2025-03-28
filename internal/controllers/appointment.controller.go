package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"med/internal/schemas"
	"net/http"
)

// CreateAppointment godoc
// @Router       /api/appointment [POST]
// @Security     ApiKeyAuth
// @Summary      Create Appointment
// @Description  Creates a new appointment
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        appointment body schemas.AppointmentPayload true "Appointment Request"
// @Success      201  {object}  schemas.ResponseSuccess
// @Failure      400  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) CreateAppointment(c *gin.Context) {
	var params schemas.AppointmentPayload
	err := c.ShouldBindJSON(&params)
	if h.handleError(c, err) {
		return
	}

	resp, err := h.services.CreateAppointment(c.Request.Context(), &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusCreated, schemas.ResponseSuccess{
		Data:   resp,
		Status: 200,
	})
}

// UpdateAppointment godoc
// @Router       /api/appointment/{id} [PUT]
// @Security     ApiKeyAuth
// @Summary      Update Appointment
// @Description  Updates an existing appointment
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        id path schemas.IDRequest true "Appointment ID"
// @Param        appointment body schemas.AppointmentPayload true "Appointment Update Request"
// @Success      200  {object}  schemas.ResponseSuccess
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) UpdateAppointment(c *gin.Context) {
	var params schemas.AppointmentPayload
	err := c.ShouldBindJSON(&params)
	if h.handleError(c, err) {
		return
	}
	id := c.Param("id")
	params.ID = id
	resp, err := h.services.UpdateAppointment(c.Request.Context(), &params)
	if h.handleError(c, err) {
		fmt.Println("=================", params)
		fmt.Println("=================", err.Error())
		return
	}

	c.JSON(http.StatusCreated, schemas.ResponseSuccess{
		Data:   resp,
		Status: 200,
	})
}

// UpdateAppointmentStatus godoc
// @Router       /api/appointment/{id}/status [PATCH]
// @Security     ApiKeyAuth
// @Summary      Update Appointment Status
// @Description  Updates the status of an existing appointment
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        id path schemas.IDRequest true "Appointment ID"
// @Param        appointment query schemas.AppointmentStatus true "Appointment Status Update Request"
// @Success      200  {object}  schemas.ResponseSuccess
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) UpdateAppointmentStatus(c *gin.Context) {
	status := c.DefaultQuery("status", "pending")
	resp, err := h.services.UpdateAppointmentStatus(c.Request.Context(), &schemas.Appointment{
		Status: schemas.AppointmentStatus(status),
		ID:     c.Param("id"),
	})
	if h.handleError(c, err) {
		return
	}
	c.JSON(http.StatusCreated, schemas.ResponseSuccess{
		Data:   resp,
		Status: 200,
	})
}

// DeleteAppointment godoc
// @Router       /api/appointment/{id} [DELETE]
// @Security     ApiKeyAuth
// @Summary      Delete Appointment
// @Description  Deletes an existing appointment by ID
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        id path string true "Appointment ID"
// @Success      200  {object}  schemas.ResponseSuccess
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) DeleteAppointment(c *gin.Context) {
	ctxWithTimeout, cancel := h.makeContext()
	defer cancel()

	params := schemas.IDRequest{
		ID: c.Param("id"),
	}

	resp, err := h.services.DeleteAppointment(ctxWithTimeout, &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusCreated, schemas.ResponseSuccess{
		Data:   resp,
		Status: 200,
	})
}

// GetUserAppointment godoc
// @Router       /api/user-appointments/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get User Appointments
// @Description  Retrieves a list of appointments for a specific user
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        filter query schemas.GetListRequestOfUserPayload false "Filter Request"
// @Success      200  {object}  schemas.ManyUserAppointment
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) GetUserAppointment(c *gin.Context) {
	ctxWithTimeout, cancel := h.makeContext()
	defer cancel()

	page := cast.ToInt(c.DefaultQuery("page", defaultPage))
	limit := cast.ToInt(c.DefaultQuery("limit", defaultLimit))
	params := schemas.GetListRequestOfUserPayload{
		Search: fmt.Sprintf("%%%s%%", c.DefaultQuery("search", "")),
		Page:   int64((page - 1) * limit),
		Limit:  int64(limit),
	}
	params.ID = c.Param("id")
	resp, err := h.services.GetUserAppointment(ctxWithTimeout, &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDoctorAppointment godoc
// @Router       /api/doctor-appointments/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get Doctor Appointments
// @Description  Retrieves a list of appointments for a specific doctor
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        id path string true "Doctor ID"
// @Param        filter query schemas.GetListRequestOfUserPayload false "Filter Request"
// @Success      200  {object}  schemas.ManyDoctorAppointment
// @Failure      400  {object}  schemas.ResponseError
// @Failure      404  {object}  schemas.ResponseError
// @Failure      500  {object}  schemas.ResponseError
func (h *Handler) GetDoctorAppointment(c *gin.Context) {
	ctxWithTimeout, cancel := h.makeContext()
	defer cancel()
	page := cast.ToInt(c.DefaultQuery("page", defaultPage))
	limit := cast.ToInt(c.DefaultQuery("limit", defaultLimit))
	params := schemas.GetListRequestOfUserPayload{
		Search: fmt.Sprintf("%%%s%%", c.DefaultQuery("search", "")),
		Page:   int64((page - 1) * limit),
		Limit:  int64(limit),
	}
	err := c.ShouldBindQuery(&params)

	params.ID = c.Param("id")
	resp, err := h.services.GetDoctorAppointment(ctxWithTimeout, &params)
	if h.handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}
