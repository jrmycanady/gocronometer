package gocronometer

type LoginResponse struct {
	Redirect string `json:"redirect"`
	Success  bool   `json:"success"`
	Error    string `json:"error"`
}
