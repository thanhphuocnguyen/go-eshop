package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
)

// RespondJSON sends a JSON response with the given status code and data
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// RespondSuccess sends a successful JSON response with data
func RespondSuccess(w http.ResponseWriter, r *http.Request, data interface{}) {
	response := dto.CreateDataResp(r.Context(), data, nil, nil)
	RespondJSON(w, http.StatusOK, response)
}

// RespondSuccessWithPagination sends a successful JSON response with data and pagination
func RespondSuccessWithPagination(w http.ResponseWriter, r *http.Request, data interface{}, pagination *dto.Pagination) {
	response := dto.CreateDataResp(r.Context(), data, pagination, nil)
	RespondJSON(w, http.StatusOK, response)
}

// RespondCreated sends a 201 Created response with data
func RespondCreated(w http.ResponseWriter, r *http.Request, data interface{}) {
	response := dto.CreateDataResp(r.Context(), data, nil, nil)
	RespondJSON(w, http.StatusCreated, response)
}

// RespondNoContent sends a 204 No Content response
func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// RespondError sends an error response with the appropriate status code
func RespondError(w http.ResponseWriter, statusCode int, errCode string, err error) {
	errorResp := dto.CreateErr(errCode, err)
	RespondJSON(w, statusCode, errorResp)
}

// RespondBadRequest sends a 400 Bad Request error response
func RespondBadRequest(w http.ResponseWriter, errCode string, err error) {
	RespondError(w, http.StatusBadRequest, errCode, err)
}

// RespondUnauthorized sends a 401 Unauthorized error response
func RespondUnauthorized(w http.ResponseWriter, errCode string, err error) {
	RespondError(w, http.StatusUnauthorized, errCode, err)
}

// RespondForbidden sends a 403 Forbidden error response
func RespondForbidden(w http.ResponseWriter, errCode string, err error) {
	RespondError(w, http.StatusForbidden, errCode, err)
}

// RespondNotFound sends a 404 Not Found error response
func RespondNotFound(w http.ResponseWriter, errCode string, err error) {
	RespondError(w, http.StatusNotFound, errCode, err)
}

// RespondInternalServerError sends a 500 Internal Server Error response
func RespondInternalServerError(w http.ResponseWriter, errCode string, err error) {
	RespondError(w, http.StatusInternalServerError, errCode, err)
}

// GetURLParam safely gets a URL parameter and validates it's not empty
func GetURLParam(r *http.Request, param string) (string, error) {
	value := chi.URLParam(r, param)
	if value == "" {
		return "", errors.New(param + " parameter is required")
	}
	return value, nil
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
