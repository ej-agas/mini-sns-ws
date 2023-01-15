package app

type ErrorResponse struct {
	Errors map[string][]string `json:"errors"`
}

func (e *ErrorResponse) Add(field, message string) {
	fieldErrors := e.Errors[field]

	e.Errors[field] = append(fieldErrors, message)
}

func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{Errors: make(map[string][]string)}
}
