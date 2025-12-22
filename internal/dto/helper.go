package dto

import (
	"context"
	"unsafe"
)

type ErrorResp struct {
	Error ApiError `json:"error"`
}

func CreateErr(code string, err error) ErrorResp {
	return ErrorResp{
		Error: ApiError{
			Code:    code,
			Details: err.Error(),
			Stack:   err,
		},
	}
}

func CreateDataResp[T any](c context.Context, data T, pagination *Pagination, err *ApiError) ApiResponse[T] {
	// Safe extraction of context values with fallbacks
	resp := ApiResponse[T]{
		Data:       &data,
		Pagination: pagination,
	}

	if err != nil {
		resp.Error = err
	}
	return resp
}

func CreatePagination(page, pageSize, total int64) *Pagination {
	return &Pagination{
		Page:            page,
		PageSize:        pageSize,
		Total:           total,
		TotalPages:      total / int64(pageSize),
		HasNextPage:     total > int64(page*pageSize),
		HasPreviousPage: page > 1,
	}
}

func IsStructEmpty(s interface{}) bool {
	return unsafe.Sizeof(s) == 0
}
