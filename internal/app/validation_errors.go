package app

type ValidationErrors struct {
	Errors map[string][]string `json:"errors"`
}

func (e *ValidationErrors) Add(field, message string) {
	fieldErrors := e.Errors[field]

	e.Errors[field] = append(fieldErrors, message)
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{Errors: make(map[string][]string)}
}
