package response

type Response struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}
