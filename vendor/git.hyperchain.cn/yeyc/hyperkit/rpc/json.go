package rpc

import "encoding/json"

const (
	JSONRPCVersion = "2.0"
)

type jsonRequest struct {
	Method  string        `json:"method"`
	Version string        `json:"jsonrpc"`
	Id      int           `json:"id, omitempty"`
	Params  []interface{} `json:"params, omitempty"`
}

type jsonResponse struct {
	Version string          `json:"jsonrpc"`
	Id      int             `json:"id, omitempty"`
	Result  json.RawMessage `json:"result, omitempty"`
	Error   jsonError       `json:"error, omitempty"`
}

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data, omitempty"`
}