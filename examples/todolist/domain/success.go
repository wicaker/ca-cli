package domain

// ResponseSuccess represent the reseponse error struct
type ResponseSuccess struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
