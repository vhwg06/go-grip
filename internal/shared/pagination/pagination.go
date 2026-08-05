package pagination

const (
	// DefaultLimit is the default page size when no limit is provided.
	DefaultLimit = 20
	// MaxLimit is the maximum upper bound for a single page request.
	MaxLimit = 100
)

// Pagination captures limit/offset request parameters.
type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// New constructs a normalized Pagination with bounded limit and offset.
func New(limit, offset int) Pagination {
	return Pagination{Limit: limit, Offset: offset}.Normalize()
}

// Normalize applies bounded defaults for public list APIs.
func (p Pagination) Normalize() Pagination {
	if p.Limit <= 0 {
		p.Limit = DefaultLimit
	}
	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	return p
}

// Page describes response pagination metadata including total record count.
type Page struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// NewPage constructs response pagination metadata.
func NewPage(limit, offset, total int) Page {
	p := Pagination{Limit: limit, Offset: offset}.Normalize()
	if total < 0 {
		total = 0
	}
	return Page{
		Limit:  p.Limit,
		Offset: p.Offset,
		Total:  total,
	}
}
