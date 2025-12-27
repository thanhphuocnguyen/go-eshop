package main

type ESStore interface {
	CreateIndex(index string, mapping interface{}) error
	IndexDocument(index string, documentID string, document interface{}) error
	QueryDocuments(index string, query interface{}) ([]interface{}, error)
	UpdateDocument(index string, documentID string, update interface{}) error
	DeleteDocument(index string, documentID string) error
}

type esStore struct {
}
