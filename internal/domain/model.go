package domain

type Model interface {
	Id() string
	String() string
}

type ModelCollection[T Model] struct {
	Data []T `json:"data"`
}

func NewModelCollection[T Model](models []T) *ModelCollection[T] {
	collection := &ModelCollection[T]{Data: make([]T, 0)}
	collection.Data = append(collection.Data, models...)

	return collection
}
