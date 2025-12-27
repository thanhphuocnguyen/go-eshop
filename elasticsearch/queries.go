package elasticsearch

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ProductQuery struct {
	Query Query `json:"query"`
}

type Query struct {
	Match Match `json:"match"`
}

type Match struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func BuildProductQuery(field, value string) (string, error) {
	if field == "" || value == "" {
		return "", fmt.Errorf("field and value must not be empty")
	}

	query := ProductQuery{
		Query: Query{
			Match: Match{
				Field: field,
				Value: value,
			},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return "", err
	}

	return string(queryJSON), nil
}

func BuildMultiFieldQuery(queries map[string]string) (string, error) {
	if len(queries) == 0 {
		return "", fmt.Errorf("queries must not be empty")
	}

	var matchQueries []Match
	for field, value := range queries {
		if field == "" || value == "" {
			return "", fmt.Errorf("field and value must not be empty")
		}
		matchQueries = append(matchQueries, Match{Field: field, Value: value})
	}

	query := ProductQuery{
		Query: Query{
			Match: Match{
				Field: strings.Join(getFields(matchQueries), ", "),
				Value: strings.Join(getValues(matchQueries), ", "),
			},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return "", err
	}

	return string(queryJSON), nil
}

func getFields(matches []Match) []string {
	fields := make([]string, len(matches))
	for i, match := range matches {
		fields[i] = match.Field
	}
	return fields
}

func getValues(matches []Match) []string {
	values := make([]string, len(matches))
	for i, match := range matches {
		values[i] = match.Value
	}
	return values
}