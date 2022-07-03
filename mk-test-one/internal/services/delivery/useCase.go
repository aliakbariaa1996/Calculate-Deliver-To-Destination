package delivery

import (
	"github.com/aliakbariaa1996/mk-test-one/config"
	loggerx "github.com/aliakbariaa1996/mk-test-one/internal/common/log"
)

type UseCase struct {
	cfg    *config.Config
	logger *loggerx.Logger
}

func NewDeliveryUseCase(cfg *config.Config, logger *loggerx.Logger) *UseCase {
	return &UseCase{
		cfg:    cfg,
		logger: logger,
	}
}

type UseService interface {
	GetDistance(sorLoc SourceLocation, deliLocs []DeliverManLocation) interface{}
	CalculateDist(sourceX float64, sourceY float64, DeliverManX float64, DeliverManY float64, c chan float64)
}
