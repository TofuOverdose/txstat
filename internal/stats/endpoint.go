package stats

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeTopExchangeDiffAddressEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		resp := &topExchangeDiffResponse{}
		addr, err := s.TopExchangeDiffAddress(ctx)
		if err != nil {
			return nil, err
		}

		resp.Address = addr
		return resp, nil
	}
}

type topExchangeDiffResponse struct {
	Address string `json:"address"`
}
