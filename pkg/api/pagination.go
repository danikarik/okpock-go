package api

// DefaultPageLimit is a default limit for pagination.
const DefaultPageLimit uint64 = 10

// NewPagingOptions returns a new instance of `PagingOptions`.
func NewPagingOptions(cursor int64, limit uint64) *PagingOptions {
	if limit == 0 {
		limit = DefaultPageLimit
	}
	return &PagingOptions{
		Cursor: cursor,
		Limit:  limit,
	}
}

// PagingOptions holds pagination fields.
type PagingOptions struct {
	Cursor int64  `json:"cursor"`
	Limit  uint64 `json:"limit"`
	Next   int64  `json:"next"`
}

// HasNext is a helper func which determines whether pager has next cursor or not.
func (p *PagingOptions) HasNext() bool { return p.Next > 0 }
