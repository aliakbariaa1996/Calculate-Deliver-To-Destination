package errorx

type Error struct {
	Code         string      `json:"code"`
	StatusCode   int         `json:"status_code"`
	Error        string      `json:"error"`
	DetailErrors interface{} `json:"detail_errors"`
}

type Success struct {
	Code    string      `json:"code"`
	Message string      `json:"error"`
	Details interface{} `json:"details"`
}
