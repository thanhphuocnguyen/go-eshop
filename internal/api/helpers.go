package api

import (
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
)

type ErrorResp struct {
	Error ApiError `json:"error"`
}

func createErr(code string, err error) ErrorResp {
	return ErrorResp{
		Error: ApiError{
			Code:    code,
			Details: err.Error(),
			Stack:   err,
		},
	}
}

func createDataResp[T any](c *gin.Context, data T, pagination *Pagination, err *ApiError) ApiResponse[T] {
	resp := ApiResponse[T]{
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

func isStructEmpty(s interface{}) bool {
	return unsafe.Sizeof(s) == 0
}
