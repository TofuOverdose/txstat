package stats

import (
	"context"
	"fmt"

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

func panicCatchingMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (res interface{}, err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				err = fmt.Errorf("panic recovered: %s", rerr)
				return
			}
		}()
		res, err = next(ctx, req)
		return
	}
}
