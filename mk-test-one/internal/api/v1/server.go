package v1

import (
	"github.com/aliakbariaa1996/mk-test-one/config"
	loggerx "github.com/aliakbariaa1996/mk-test-one/internal/common/log"
	"github.com/aliakbariaa1996/mk-test-one/internal/services/delivery"

	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo

	logger  *loggerx.Logger
	ss      *ServiceStorage
	cfg     *config.Config
	handler Handler
}

type ServiceStorage struct {
	deliveryService delivery.UseService
}

type Handler struct {
	logger *loggerx.Logger
	cfg    *config.Config
}

func NewServer(router *echo.Echo, cfg *config.Config, logger *loggerx.Logger) (*Server, error) {
	var err error
	s := &Server{
		Echo:   router,
		cfg:    cfg,
		logger: logger,
	}
	s.ss = NewServiceStorage(cfg, logger)
	s.handler = Handler{logger: logger, cfg: cfg}

	// routes init
	s.initRoutes()
	return s, err
}

func NewServiceStorage(cfg *config.Config, logger *loggerx.Logger) *ServiceStorage {
	return &ServiceStorage{
		deliveryService: delivery.NewDeliveryUseCase(cfg, logger),
	}
}
