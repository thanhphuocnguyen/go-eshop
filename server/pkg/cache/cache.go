package cache

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (string, error)
	Delete(key string) error
}
