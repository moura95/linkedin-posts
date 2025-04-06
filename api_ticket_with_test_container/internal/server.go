package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github/moura95/ticket-api/config"
	"github/moura95/ticket-api/docs"
	"github/moura95/ticket-api/internal/api"
	"github/moura95/ticket-api/internal/middleware"
	"github/moura95/ticket-api/internal/repository"
	"github/moura95/ticket-api/internal/service"

	"go.uber.org/zap"
)

type Server struct {
	store  *repository.Querier
	router *gin.Engine
	config *config.Config
	logger *zap.SugaredLogger
}

//      @title                  ticket Api
//      @version                1.0
//      @description    Api Tickets
//      @termsOfService http://swagger.io/terms/

//      @license.name   Apache 2.0
//      @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host           localhost:8080
// @BasePath       /api/v1
func NewServer(cfg config.Config, store repository.Querier, log *zap.SugaredLogger) *Server {

	server := &Server{
		store:  &store,
		config: &cfg,
		logger: log,
	}
	var router *gin.Engine

	router = gin.Default()

	docs.SwaggerInfo.BasePath = ""

	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Middleware Rate Limiter
	router.Use(middleware.RateLimitMiddleware())

	// Middleware Cors
	router.Use(middleware.CORSMiddleware())

	// Init all Routers
	createRoutesV1(&store, server.config, router, log)

	server.router = router
	return server
}

func createRoutesV1(store *repository.Querier, cfg *config.Config, router *gin.Engine, log *zap.SugaredLogger) {
	routes := router.Group("/")
	// Instance Ticket Service
	ticketService := service.NewTicketService(*store, *cfg, log)
	api.NewTicketRouter(*ticketService, log).SetupTicketRoute(routes)

	// Instance Category Service
	categoryService := service.NewCategoryService(*store, *cfg, log)
	api.NewCategoryRouter(*categoryService, log).SetupCategoryRoute(routes)

	// User Route
	api.NewUserRouter(log).SetupUserRoute(routes)

}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func RunGinServer(cfg config.Config, store repository.Querier, log *zap.SugaredLogger) {
	server := NewServer(cfg, store, log)

	_ = server.Start(cfg.HTTPServerAddress)
}
