package elasticsearch

type Config struct {
	ElasticsearchURL string `json:"elasticsearch_url"`
	IndexName        string `json:"index_name"`
	Username         string `json:"username,omitempty"`
	Password         string `json:"password,omitempty"`
}