package eth

import "encoding/json"

type JsonRpcRequest struct {
	Version string `json:"jsonrpc"`
	Method  string `json:"method"`
	Id      string `json:"id"`
	Params  interface{}
}

type JsonRpcResponse struct {
	Version string          `json:"jsonrpc"`
	Id      string          `json:"id"`
	Result  json.RawMessage `json:"result"`
}

func decodeBlockNumberResult(b []byte) (interface{}, error) {
	var res string
	err := json.Unmarshal(b, &res)
	return res, err
}

func decodeGetBlockByNumberResult(b []byte) (interface{}, error) {
	var res Block
	err := json.Unmarshal(b, &res)
	return res, err
}
