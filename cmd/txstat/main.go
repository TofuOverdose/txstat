package main

import (
	"log"
	nethttp "net/http"
	"os"

	"github.com/tofuoverdose/txstat/internal/txstat/domain/stats"
	"github.com/tofuoverdose/txstat/internal/txstat/fetcher"
	"github.com/tofuoverdose/txstat/internal/txstat/ports/http"
	"github.com/tofuoverdose/txstat/pkg/getblock/eth"
)

func main() {
	ethClientUrl, ok := os.LookupEnv("TXSTAT_GETBLOCKIO_ENDPOINT")
	if !ok {
		log.Println("Fatal error: missing env variable TXSTAT_GETBLOCKIO_ENDPOINT")
		os.Exit(1)
	}
	ethClientToken, ok := os.LookupEnv("TXSTAT_GETBLOCKIO_TOKEN")
	if !ok {
		log.Println("Fatal error: missing env variable TXSTAT_GETBLOCKIO_TOKEN")
		os.Exit(1)
	}

	ethClient := eth.Client{
		HttpClient: nethttp.Client{},
		Url:        ethClientUrl,
		Token:      ethClientToken,
	}

	f := &fetcher.Fetcher{Client: ethClient}

	statsService := stats.NewService(f)

	httpServerPort, ok := os.LookupEnv("TXSTAT_HTTP_SERVER_ENDPOINT")
	if !ok {
		log.Println("Fatal error: missing env variable TXSTAT_HTTP_SERVER_ENDPOINT")
		os.Exit(1)
	}

	httpServerConfig := http.ServerConfig{
		Port: httpServerPort,
	}
	httpServer, err := http.NewServer(httpServerConfig, statsService)
	if err != nil {
		panic(err)
	}

	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
