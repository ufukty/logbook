package response

type UserCreate struct {
	Status       bool   `json:"status"`
	ErrorMessage string `json:"error_message"`
}

type UserDelete struct {
	Status       bool   `json:"status"`
	ErrorMessage string `json:"error_message"`
}
