package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"med/internal/configs"
	repository "med/internal/db"
	"med/internal/services"
	"med/pkg/email"
	"med/pkg/logger"
)

type Server interface {
	Run()
	Stop()
}

type Handler struct {
	Engine    *gin.Engine
	services  *services.Service
	log       logger.ILogger
	cnf       configs.Config
	scheduler *gocron.Scheduler
}

func NewHandler(engine *gin.Engine, services *services.Service, log logger.ILogger, cnf configs.Config) *Handler {
	return &Handler{
		Engine:    engine,
		services:  services,
		log:       log,
		cnf:       cnf,
		scheduler: gocron.NewScheduler(time.UTC),
	}
}

func NewServer(cfg configs.Config) Server {
	loggerLevel := logger.LevelDebug
	switch cfg.Environment {
	case configs.DebugMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case configs.TestMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}
	log := logger.NewLogger("med.uz", loggerLevel)
	defer logger.Cleanup(log)

	engine := gin.Default()
	defaultConfig := cors.DefaultConfig()
	defaultConfig.AllowCredentials = true
	defaultConfig.AllowOrigins = []string{"*"}
	defaultConfig.AllowHeaders = []string{
		"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token",
		"Authorization", "accept", "origin", "Cache-Control", "X-Requested-With",
	}
	defaultConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	engine.Use(cors.New(defaultConfig))

	dbTx, err := repository.New(context.Background(), cfg, log)
	if err != nil {
		log.Error("Failed to initialize database repository", logger.Error(err))
		panic(err)
	}
	eml := email.NewEmail(&cfg)
	service := services.NewService(dbTx, dbTx.Pool(), log, eml)

	handler := NewHandler(engine, service, log, cfg)
	handler.scheduler.Every(5).Minute().Do(handler.services.ProcessEmailNotifications)
	handler.scheduler.StartAsync()

	SetUpApi(handler)

	return handler
}

// @title MED App
// @description This API contains the source for the med.uz app
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath

// Run initializes http server
func (h *Handler) Run() {
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL(fmt.Sprintf(
			"%s:%d/swagger/docs.json",
			h.cnf.HTTPHost,
			h.cnf.HTTPPort,
		)),
		ginSwagger.DefaultModelsExpandDepth(-1),
	)

	h.log.Info("server is running: ", logger.Any("address", fmt.Sprintf("%s:%d", h.cnf.HTTPHost, h.cnf.HTTPPort)))
	h.log.Info("swagger: ", logger.Any("url", fmt.Sprintf("http://%s:%d/swagger/index.html", h.cnf.HTTPHost, h.cnf.HTTPPort)))

	if err := h.Engine.Run(fmt.Sprintf(":%d", h.cnf.HTTPPort)); err != nil {
		h.log.Error("failed to run server", logger.Error(err))
	}
}

func (h *Handler) Stop() {
	h.scheduler.Stop() // Scheduler'ni to‘xtatish
	h.log.Info("shutting down")
}
