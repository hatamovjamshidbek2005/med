package handlers

import (
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"med/internal/auth"
	"med/pkg/logger"
	ratelimiter "med/pkg/ratelimter"
)

func SetUpApi(h *Handler) {
	// parse qilib oladi yani configs.yaml filega ratelimter ning bir api sorovlar keltriilgan
	config, err := ratelimiter.ParseYamlFile("./configs/config.yaml")
	if err != nil {
		_ = fmt.Errorf("failed to parse rate limiter config: %w", logger.Error(err))
	}

	rl, err := ratelimiter.NewRateLimiter(config)
	if err != nil {
		_ = fmt.Errorf("failed to create rate limiter: %w", err)
	}

	v1 := h.Engine.Group("/api")
	v1.Use(rl.GinMiddleware())
	{
		// Auth routerlari
		v1.POST("/auth/register", h.RegisterUser)
		v1.POST("/auth/login", h.LoginUser)
		v1.PUT("/auth/password", h.UpdatePassword)

	}
	v1.Use(auth.AuthMiddleware())
	// bu api larga casbin qoyip ketsa boladi qisqa tarzida men qoyip o'tirmadim deyarli kichik bolgani sababli
	// proyect assosiy maqsadi code bolsa kerak b
	// bu proyect ratelimter ishlatilgan yani bir nechta request bir nechta soniyalarda toxtovsiz kelsa request ni blocklaydi

	{
		// Appointment
		v1.POST("/appointment", h.CreateAppointment)
		v1.PUT("/appointment/:id", h.UpdateAppointment)
		v1.PATCH("/appointment/:id/status", h.UpdateAppointmentStatus)
		v1.DELETE("/appointment/:id", h.DeleteAppointment)
		v1.GET("/user-appointments/:id", h.GetUserAppointment)
		v1.GET("/doctor-appointments/:id", h.GetDoctorAppointment)

		// Doctor routerlari
		v1.POST("/doctor", h.CreateDoctor)
		v1.PUT("/doctor/:id", h.UpdateDoctor)
		v1.GET("/doctor/:id", h.GetDoctor)
		v1.DELETE("/doctor/:id", h.DeleteDoctor)
		v1.GET("/doctors/list", h.GetAllDoctors)

		//User routerlari
		v1.PUT("/user/:id", h.UpdateUser)
		v1.DELETE("/user/:id", h.DeleteUser)
		v1.GET("/user/:id", h.GetUser)
		v1.GET("/users/list", h.GetAllUsers)
	}
	h.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
