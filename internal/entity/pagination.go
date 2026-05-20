package entity

// Pagination captures limit/offset request metadata.
type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Normalize applies bounded defaults for public list APIs.
func (p Pagination) Normalize() Pagination {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	return p
}

// Page describes paginated response metadata.
type Page struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
