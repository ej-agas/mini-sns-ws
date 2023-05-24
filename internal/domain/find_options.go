package domain

type SortOrder int

const (
	SortAscending  SortOrder = 1
	SortDescending SortOrder = -1
)

type FindOptions struct {
	Sort  map[string]SortOrder
	Limit int64
}

func NewFindOptions() *FindOptions {
	return &FindOptions{}
}

func (opts *FindOptions) AddSortBy(field string, order SortOrder) *FindOptions {
	opts.Sort[field] = order

	return opts
}

func (opts *FindOptions) SetLimit(limit int64) *FindOptions {
	opts.Limit = limit

	return opts
}
