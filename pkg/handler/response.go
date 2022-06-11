package handler

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

const (
	StatusFail    = "fail"
	StatusSuccess = "success"
)
