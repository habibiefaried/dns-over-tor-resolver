package cachehandler

type CacheHandler interface {
	Init() error
	Put(string, string, string, int64) error // set Key value
	Get(string) ([]string, error)            // get value
	Close() error
	CleanUp() (int64, error) // cleanup all records aged than N seconds
}
