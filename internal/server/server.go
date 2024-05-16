package server

import (
	docs "employer.dev/docs"
	"employer.dev/internal/api"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	config Config
	api    *api.EmployerAPI
	logger *logrus.Logger
}

var router *gin.Engine

func init() {
	router = gin.Default()
}

func NewServer(config Config, api *api.EmployerAPI, logger *logrus.Logger) *Server {
	return &Server{
		config: config,
		api:    api,
		logger: logger,
	}
}

func (s *Server) Start() error {
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%s/swagger/doc.json", s.config.host, s.config.port)), // The url pointing to API definition
	)))
	s.api.RegisterRoutes(router)
	
	return router.Run(s.config.host + ":" + s.config.port)
}
