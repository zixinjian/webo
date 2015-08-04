package svc

type Params map[string]interface{}

type Results struct {
	Status  string      `json:"status"`
	Results interface{} `json:"results"`
}

var (
	MaxLimit         int64 = 10000
	DefaultRowsLimit int64 = 50
)
