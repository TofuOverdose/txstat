package stats

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	gormux "github.com/gorilla/mux"
)

func RegisterHttpServer(mux *gormux.Router, service Service) *gormux.Router {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeErrorFunc),
	}

	{
		h := kithttp.NewServer(
			panicCatchingMiddleware(makeTopExchangeDiffAddressEndpoint(service)),
			kithttp.NopRequestDecoder,
			genericJsonResponseEncoder,
			opts...,
		)
		mux.Handle("/exchange/top", h).Methods("GET")
	}

	return mux
}

func genericJsonResponseEncoder(_ context.Context, w http.ResponseWriter, res interface{}) error {
	resp := response{}
	if failer, ok := res.(endpoint.Failer); ok {
		if err := failer.Failed(); err != nil {
			resp.Error = err.Error()
			w.WriteHeader(400)
		}
	} else {
		resp.Result = res
		w.WriteHeader(200)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, _ = w.Write(b)
	return nil
}

func encodeErrorFunc(_ context.Context, _ error, w http.ResponseWriter) {
	resp := response{
		Error: "internal error",
	}
	b, _ := json.Marshal(resp)
	w.WriteHeader(500)
	_, _ = w.Write(b)
}

type response struct {
	Result interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}
