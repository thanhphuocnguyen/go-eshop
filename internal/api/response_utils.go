package api

import (
	"encoding/json"
	"net/http"

	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
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
func RespondSuccessWithError(w http.ResponseWriter, r *http.Request, data interface{}, err *dto.ApiError) {
	response := dto.CreateDataResp(r.Context(), data, nil, err)
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
