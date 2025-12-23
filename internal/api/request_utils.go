package api

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
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

// GetFormData parses form data (multipart/form-data or application/x-www-form-urlencoded)
// into a struct pointed to by dest. The dest parameter must be a pointer to a struct.
// Struct fields should have `form` tags to specify the form field names.
func (s *Server) GetFormData(r *http.Request, dest interface{}) error {
	// Check if dest is a pointer
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer")
	}

	destElem := destValue.Elem()
	if destElem.Kind() != reflect.Struct {
		return errors.New("dest must be a pointer to a struct")
	}

	// Parse the form data
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		// Fallback to parsing regular form if multipart fails
		err = r.ParseForm()
		if err != nil {
			return errors.New("failed to parse form data: " + err.Error())
		}
	}

	destType := destElem.Type()
	// Iterate through struct fields
	for i := 0; i < destElem.NumField(); i++ {
		field := destElem.Field(i)
		fieldType := destType.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Get the form tag or use field name
		formTag := fieldType.Tag.Get("form")
		if formTag == "" {
			formTag = strings.ToLower(fieldType.Name)
		}

		// Skip if tag is "-"
		if formTag == "-" {
			continue
		}

		// Handle file uploads for multipart.FileHeader fields
		if field.Type() == reflect.TypeOf((*multipart.FileHeader)(nil)) {
			if r.MultipartForm != nil && r.MultipartForm.File != nil {
				if files, exists := r.MultipartForm.File[formTag]; exists && len(files) > 0 {
					field.Set(reflect.ValueOf(files[0]))
				}
			}
			continue
		}

		// Get the form value for non-file fields
		formValue := r.FormValue(formTag)
		if formValue == "" {
			continue
		}

		// Set the field value based on its type
		err = setFieldValue(field, formValue)
		if err != nil {
			return errors.New("error setting field " + fieldType.Name + ": " + err.Error())
		}
	}

	// Validate the struct using the validator
	err = s.validator.Struct(dest)
	if err != nil {
		return err
	}

	return nil
}

// setFieldValue sets the value of a struct field based on its type
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	case reflect.Ptr:
		// Handle pointer fields
		if field.Type().Elem().Kind() == reflect.String {
			field.Set(reflect.ValueOf(&value))
		} else {
			// For other pointer types, create a new value and set it
			newVal := reflect.New(field.Type().Elem())
			err := setFieldValue(newVal.Elem(), value)
			if err != nil {
				return err
			}
			field.Set(newVal)
		}
	default:
		return errors.New("unsupported field type: " + field.Kind().String())
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

// GetUserClaimsFromContext safely extracts JWT claims from request context
func GetUserClaimsFromContext(r *http.Request) (map[string]interface{}, error) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// GetUserIDFromClaims safely extracts user ID from JWT claims
func GetUserIDFromClaims(claims map[string]interface{}) (uuid.UUID, error) {
	userIdValue, exists := claims["userId"]
	if !exists {
		return uuid.Nil, errors.New("userId not found in token claims")
	}

	// Try to handle different possible types
	switch v := userIdValue.(type) {
	case uuid.UUID:
		return v, nil
	case string:
		return uuid.Parse(v)
	default:
		return uuid.Nil, errors.New("userId in token claims has invalid type")
	}
}

// GetRoleCodeFromClaims safely extracts role code from JWT claims
func GetRoleCodeFromClaims(claims map[string]interface{}) (string, error) {
	roleCodeValue, exists := claims["roleCode"]
	if !exists {
		return "", errors.New("roleCode not found in token claims")
	}

	roleCode, ok := roleCodeValue.(string)
	if !ok {
		return "", errors.New("roleCode in token claims has invalid type")
	}
	return roleCode, nil
}
