package stats

import (
	"context"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
)

type loggingService struct {
	service     Service
	infoLogger  log.Logger
	errorLogger log.Logger
}

func NewLoggingService(service Service, logger log.Logger) Service {
	return loggingService{
		service:     service,
		infoLogger:  level.Info(logger),
		errorLogger: level.Error(logger),
	}
}

func (s loggingService) TopExchangeDiffAddress(ctx context.Context) (res string, err error) {
	start := time.Now()
	s.infoLogger.Log("msg", "Start TopExchangeDiffAddress")

	res, err = s.service.TopExchangeDiffAddress(ctx)

	if err != nil {
		s.errorLogger.Log("msg", "Finish TopExchangeDiffAddress with error",
			"elapsed", time.Since(start),
			"error", err.Error())
		return
	}
	s.infoLogger.Log("msg", "Finish TopExchangeDiffAddress", "elapsed", time.Since(start))
	return
}
