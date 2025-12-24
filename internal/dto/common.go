package dto

// Response types - unchanged
type ApiResponse[T any] struct {
	Message    string      `json:"message,omitempty"`
	Data       *T          `json:"data,omitempty,omitzero"`
	Error      *ApiError   `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
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
