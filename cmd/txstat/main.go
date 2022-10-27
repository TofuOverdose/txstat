package main

import (
	"fmt"
	"net"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"

	"github.com/tofuoverdose/txstat/internal/fetcher"
	"github.com/tofuoverdose/txstat/internal/stats"
	"github.com/tofuoverdose/txstat/pkg/getblock/eth"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))

	router := mux.NewRouter()
	router.Use(requestLoggingMiddleware(logger))

	var statsService stats.Service
	{
		f := &fetcher.Fetcher{
			Client: eth.Client{
				HttpClient: nethttp.Client{},
				Url:        mustGetEnv("GETBLOCKIO_ETH_URL"),
				Token:      mustGetEnv("GETBLOCKIO_ETH_TOKEN"),
			},
		}

		statsService = stats.NewService(f)
		statsService = stats.NewLoggingService(statsService, logger)
		stats.RegisterHttpServer(router.PathPrefix("/stats").Subrouter(), statsService)
	}

	errs := make(chan error, 2)
	go func() {
		addr := net.JoinHostPort("", mustGetEnv("HTTP_SERVER_PORT"))
		level.Info(logger).Log("msg", "starting http server", "addr", addr)
		errs <- nethttp.ListenAndServe(addr, router)
	}()
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-sig)
	}()

	logger.Log("process exited", <-errs)
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("env variable %s is missing\n", key))
	}
	return v
}

func requestLoggingMiddleware(logger log.Logger) mux.MiddlewareFunc {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			now := time.Now()
			logger = log.With(logger, "uri", r.RequestURI)
			logger.Log("msg", "received http request")
			defer func(start time.Time) {
				logger.Log("msg", "sent http response", "elapsed", time.Since(start))
			}(now)
			next.ServeHTTP(w, r)
		})
	}
}
