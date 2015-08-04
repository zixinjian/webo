package controllers

type JsonResult struct {
	Ret    string      `json:"ret"`
	Result interface{} `json:"result"`
}

type TableResult struct {
	Status string      `json:status`
	Total  int64       `json:"total"`
	Rows   interface{} `json:"rows"`
}
