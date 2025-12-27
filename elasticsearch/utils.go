package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// ToReader transforms various data types to io.Reader for Elasticsearch operations.
// Supported types:
// - string: returns a strings.Reader
// - []byte: returns a bytes.Reader
// - io.Reader: returns as-is
// - struct/map/slice: marshals to JSON and returns bytes.Reader
// - nil: returns an empty bytes.Reader
func ToReader(data interface{}) (io.Reader, error) {
	if data == nil {
		return bytes.NewReader([]byte{}), nil
	}

	switch v := data.(type) {
	case string:
		return strings.NewReader(v), nil
	case []byte:
		return bytes.NewReader(v), nil
	case io.Reader:
		return v, nil
	default:
		// For structs, maps, slices, and other types, marshal to JSON
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data to JSON: %w", err)
		}
		return bytes.NewReader(jsonData), nil
	}
}

// ToReaderWithValidation transforms data to io.Reader with additional validation.
// It ensures that the data is not empty and can be properly serialized.
func ToReaderWithValidation(data interface{}) (io.Reader, error) {
	if data == nil {
		return nil, fmt.Errorf("data cannot be nil")
	}

	// Check if it's an empty struct, slice, or map
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Slice, reflect.Map:
		if v.Len() == 0 {
			return nil, fmt.Errorf("data cannot be empty")
		}
	case reflect.String:
		if v.String() == "" {
			return nil, fmt.Errorf("data cannot be empty string")
		}
	case reflect.Ptr:
		if v.IsNil() {
			return nil, fmt.Errorf("data cannot be nil pointer")
		}
	}

	return ToReader(data)
}

// MustToReader is like ToReader but panics on error.
// Use this only when you're certain the data conversion will succeed.
func MustToReader(data interface{}) io.Reader {
	reader, err := ToReader(data)
	if err != nil {
		panic(fmt.Sprintf("failed to convert data to io.Reader: %v", err))
	}
	return reader
}

// ToReaderPretty transforms data to io.Reader with pretty-printed JSON.
// Useful for debugging or when you need formatted JSON output.
func ToReaderPretty(data interface{}) (io.Reader, error) {
	if data == nil {
		return bytes.NewReader([]byte{}), nil
	}

	switch v := data.(type) {
	case string:
		return strings.NewReader(v), nil
	case []byte:
		return bytes.NewReader(v), nil
	case io.Reader:
		return v, nil
	default:
		// For structs, maps, slices, marshal to pretty JSON
		jsonData, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data to pretty JSON: %w", err)
		}
		return bytes.NewReader(jsonData), nil
	}
}
