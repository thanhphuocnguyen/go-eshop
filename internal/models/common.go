package models

type URISlugParam struct {
	Slug string `uri:"slug" binding:"required"`
}

type UriIDParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}
type PublicIDParam struct {
	PublicID string `uri:"publicId" binding:"required"`
}

// Meta info about the request
type MetaInfo struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId"`
	Path      string `json:"path"`
	Method    string `json:"method"`
}

// Response types - unchanged
type ApiResponse[T any] struct {
	Message    string      `json:"message"`
	Data       *T          `json:"data,omitempty,omitzero"`
	Error      *ApiError   `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Meta       *MetaInfo   `json:"meta"`
}

// Error structure for detailed errors
type ApiError struct {
	Code    string `json:"code"`
	Details string `json:"details"`
	Stack   error  `json:"stack,omitempty"` // Hide in production
}

// Pagination info (for paginated endpoints)
type Pagination struct {
	Total           int64 `json:"total"`
	Page            int64 `json:"page"`
	PageSize        int64 `json:"pageSize"`
	TotalPages      int64 `json:"totalPages"`
	HasNextPage     bool  `json:"hasNextPage"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
}
