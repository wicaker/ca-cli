package domain

// ResponseSuccess represent the reseponse success struct
type ResponseSuccess struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
