package model

import "time"

// FilterCriteria defines the filtering and pagination options for querying audit logs.
type FilterCriteria struct {
	ServiceName   string
	Username      string
	Action        string
	ResourceType  string
	ResourceID    string
	TimestampFrom time.Time
	TimestampTo   time.Time
	Page          int
	PageSize      int
}

// GetLimit returns the page size for pagination.
func (f FilterCriteria) GetLimit() int {
	if f.PageSize > 0 {
		return f.PageSize
	}
	return 20
}

// GetOffset returns the offset for pagination based on page and page size.
func (f FilterCriteria) GetOffset() int {
	if f.Page > 0 && f.PageSize > 0 {
		return (f.Page - 1) * f.PageSize
	}
	return 0
}
