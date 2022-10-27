package http

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type GetGreatestBalanceDiffResponse struct {
	Address string `json:"address"`
}
