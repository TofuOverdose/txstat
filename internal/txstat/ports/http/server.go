package http

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/tofuoverdose/txstat/internal/txstat/domain/stats"
)

type Server struct {
	*http.Server
	service *stats.Service
}

const (
	greatestBalanceDiffPath = "/stats/greatestBalanceDiff"
)

type ServerConfig struct {
	Port string
}

func NewServer(cfg ServerConfig, service *stats.Service) (*Server, error) {
	s := Server{
		service: service,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(greatestBalanceDiffPath, s.handleGetGreatestBalanceDiff)

	s.Server = &http.Server{
		Addr:         net.JoinHostPort("", cfg.Port),
		Handler:      mux,
		WriteTimeout: 5 * time.Minute,
	}

	return &s, nil
}

func (s *Server) handleGetGreatestBalanceDiff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}

	resp := Response{}
	addr, err := s.service.GetAddressWithGreatestExchangeDiff(r.Context())
	if err != nil {
		if errors.Is(err, stats.ErrEmptyBlockChain) {
			w.WriteHeader(404)
			resp.Error = "no data in block chain"
		} else {
			log.Printf("failed to GetAddressWithGreatestExchangeDiff: err=%s", err.Error())
			w.WriteHeader(500)
			resp.Error = "internal error"
		}
		b, _ := json.Marshal(resp)
		_, _ = w.Write(b)
		return
	}

	respData := GetGreatestBalanceDiffResponse{
		Address: addr,
	}
	resp.Data = respData

	b, _ := json.Marshal(resp)
	_, _ = w.Write(b)
}
