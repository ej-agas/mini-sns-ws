package domain

type Models interface {
	User | Post
}

type ModelCollection[T Models] struct {
	Data []T `json:"data"`
}

func NewModelCollection[T Models](models []T) *ModelCollection[T] {
	collection := &ModelCollection[T]{Data: make([]T, 0)}
	collection.Data = append(collection.Data, models...)

	return collection
}
