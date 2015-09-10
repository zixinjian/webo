package t

type ItemMap map[string]interface{}

type Params map[string]interface{}
type LimitParams map[string]int64

type Results struct {
	Status  string      `json:"status"`
	Results interface{} `json:"results"`
}

const (
	LimitDefault int64 = 100
	LimitMax     int64 = 10000
)
