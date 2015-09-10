package rpc

import (
	"encoding/json"
)

type JsonRequest struct {
	Id     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

type JsonResponse struct {
	Id     *json.RawMessage `json:"id"`
	Result interface{}      `json:"result"`
	Error  interface{}      `json:"error"`
}

type TableResult struct {
	Total int         `json:"total"`
	Rows  interface{} `json:"rows"`
}

type JsonResult struct {
	Ret    string      `json:"ret"`
	Result interface{} `json:"result"`
}
