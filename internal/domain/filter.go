package domain

type Filter = map[string]interface{}

func NewFilter() Filter {
	return make(map[string]interface{})
}
