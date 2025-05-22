package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

func createErrorResponse[T any](code string, msg string, err error) ApiResponse[T] {
	return ApiResponse[T]{
		Success: false,
		Data:    nil,
		Error: &ApiError{
			Code:    code,
			Details: msg,
			Stack:   err.Error(),
		},
	}
}

func createSuccessResponse[T any](c *gin.Context, data T, message string, pagination *Pagination, err *ApiError) ApiResponse[T] {
	resp := ApiResponse[T]{
		Success:    true,
		Message:    message,
		Data:       &data,
		Pagination: pagination,

		Meta: &MetaInfo{
			Timestamp: time.Now().Format(time.RFC3339),
			RequestID: c.GetString("RequestID"),
			Path:      c.FullPath(),
			Method:    c.Request.Method,
		},
	}
	if err != nil {
		resp.Error = err
	}
	return resp
}

func createPagination(page, pageSize, total int64) *Pagination {
	return &Pagination{
		Page:            page,
		PageSize:        pageSize,
		Total:           total,
		TotalPages:      total / int64(pageSize),
		HasNextPage:     total > int64(page*pageSize),
		HasPreviousPage: page > 1,
	}
}
