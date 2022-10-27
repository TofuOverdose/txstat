package eth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	HttpClient http.Client
	Url        string
	Token      string
}

func (c *Client) BlockNumber(ctx context.Context) (string, error) {
	v, err := c.do(ctx, "eth_blockNumber", nil, decodeBlockNumberResult)
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

func (c *Client) GetBlockByNumber(ctx context.Context, blockNumber string, withTransactionData bool) (*Block, error) {
	params := []interface{}{blockNumber, withTransactionData}
	v, err := c.do(ctx, "eth_getBlockByNumber", params, decodeGetBlockByNumberResult)
	if err != nil {
		return nil, err
	}
	res := v.(Block)
	return &res, nil
}

func (c *Client) do(ctx context.Context, method string, params []interface{}, respDecodeFunc func([]byte) (interface{}, error)) (interface{}, error) {
	jrreq := JsonRpcRequest{
		Version: "2.0",
		Method:  method,
		Params:  params,
	}
	b, err := json.Marshal(jrreq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal jsonrpc request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.Token)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client returned error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response code is %s", resp.Status)
	}

	jrresp := new(JsonRpcResponse)
	if err := json.NewDecoder(resp.Body).Decode(jrresp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal jsonrpc response: %w", err)
	}

	result, err := respDecodeFunc(jrresp.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal result data: %w", err)
	}

	return result, nil
}
