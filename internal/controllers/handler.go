package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"med/pkg/logger"
)

var (
	defaultPage  = "1"
	defaultLimit = "10"
)

var (
	errorUsernameUnique = errors.New("ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)")
	errorIDNotFound     = errors.New("no rows in result set")
	errorStorage        = errors.New("storage error")
	errorBadRequest     = errors.New("bad request")
	errorInternal       = errors.New("internal error")
	errorAccessDenied   = errors.New("access denied")
	errorInvalidID      = errors.New("ERROR: invalid input syntax for type uuid (SQLSTATE 22P02)")
)

var (
	msgUsernameTaken = "Bu foydalanuvchi nomi allaqachon band qilingan"
	msgIDNotFound    = "Bunday ID bilan ma'lumot topilmadi"
	msgStorageError  = "Ma'lumotlarni saqlashda xatolik yuz berdi"
	msgBadRequest    = "So‘rov noto‘g‘ri kiritildi"
	msgInternalError = "Serverda xatolik yuz berdi, qayta urinib ko‘ring"
	msgAccessDenied  = "Ruxsat yo‘q"
	msgInvalidID     = "ID noto‘g‘ri kiritildi"
)

func (h *Handler) handleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	h.log.Error("Xato yuz berdi", logger.Error(err))

	switch err.Error() {
	case errorUsernameUnique.Error():
		c.JSON(http.StatusBadRequest, gin.H{"error": msgUsernameTaken})
	case errorIDNotFound.Error(), sql.ErrNoRows.Error():
		c.JSON(http.StatusNotFound, gin.H{"error": msgIDNotFound})
	case errorStorage.Error():
		c.JSON(http.StatusInternalServerError, gin.H{"error": msgStorageError})
	case errorBadRequest.Error():
		c.JSON(http.StatusBadRequest, gin.H{"error": msgBadRequest})
	case errorAccessDenied.Error():
		c.JSON(http.StatusForbidden, gin.H{"error": msgAccessDenied})
	case errorInvalidID.Error():
		c.JSON(http.StatusBadRequest, gin.H{"error": msgInvalidID})
	case errorInternal.Error():
		c.JSON(http.StatusInternalServerError, gin.H{"error": msgInternalError})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %v", msgInternalError, err)})
	}

	return true
}

func (h *Handler) makeContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Hour)
	return ctx, cancel
}
