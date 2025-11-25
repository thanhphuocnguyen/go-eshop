package models

type ImageModel struct {
	ID          string   `json:"id"`
	ExternalID  string   `json:"externalId"`
	Url         string   `json:"url"`
	MimeType    string   `json:"mimeType,omitempty"`
	FileSize    int64    `json:"fileSize,omitzero"`
	Assignments []string `json:"assignments,omitempty"`
}
