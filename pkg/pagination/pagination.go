package pagination

// Result represents offset pagination metadata and items.
type Result[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// Options represents offset pagination parameters.
type Options struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

// ApplyDefaults normalizes user-provided values.
func (p *Options) ApplyDefaults() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 || p.PageSize > 100 {
		p.PageSize = 10
	}
}

// Offset returns the SQL offset for the current page.
func (p Options) Offset() int {
	return (p.Page - 1) * p.PageSize
}
