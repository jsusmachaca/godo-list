package response

import "github.com/jsusmachaca/godo/pkg/model"

type Response struct {
	Success bool       `json:"success"`
	Data    model.Task `json:"data"`
}
