package api

import (
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
)

type ErrorResp struct {
	Error dto.ApiError `json:"error"`
}

func createErr(code string, err error) ErrorResp {
	return ErrorResp{
		Error: dto.ApiError{
			Code:    code,
			Details: err.Error(),
			Stack:   err,
		},
	}
}

func createDataResp[T any](c *gin.Context, data T, pagination *dto.Pagination, err *dto.ApiError) dto.ApiResponse[T] {
	resp := dto.ApiResponse[T]{
		Data:       &data,
		Pagination: pagination,
		Meta: &dto.MetaInfo{
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

func createPagination(page, pageSize, total int64) *dto.Pagination {
	return &dto.Pagination{
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
