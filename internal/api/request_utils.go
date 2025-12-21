package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
)

func (s *Server) GetRequestBody(r *http.Request, dest interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dest)
	if err != nil {
		return err
	}
	err = s.validator.Struct(dest)
	if err != nil {
		return err
	}
	return nil
}

func GetUrlParam(r *http.Request, key string) (string, error) {
	value := chi.URLParam(r, key)
	if value == "" {
		return "", errors.New("url parameter " + key + " is required")
	}
	return value, nil
}

// GetPaginationQuery parses URL query parameters into a PaginationQuery struct
func GetPaginationQuery(r *http.Request) models.PaginationQuery {
	query := models.PaginationQuery{
		Page:     1,  // default value
		PageSize: 20, // default value
		Search:   nil,
	}

	// Parse page parameter
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.ParseInt(pageStr, 10, 64); err == nil && page >= 1 {
			query.Page = page
		}
	}

	// Parse pageSize parameter
	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		if pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64); err == nil && pageSize >= 1 && pageSize <= 100 {
			query.PageSize = pageSize
		}
	}

	// Parse search parameter
	if searchStr := r.URL.Query().Get("search"); searchStr != "" && len(searchStr) <= 1000 {
		query.Search = &searchStr
	}

	return query
}

// ParsePaginationQuery parses standard pagination query parameters
func ParsePaginationQuery(r *http.Request) models.PaginationQuery {
	var query models.PaginationQuery
	query.Page = 1
	query.PageSize = 10

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			query.Page = int64(p)
		}
	}
	if pageSize := r.URL.Query().Get("pageSize"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			query.PageSize = int64(ps)
		}
	}
	return query
}
