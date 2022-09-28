// Stand-in for some kind of key-value distributed storage, like etcd
package dist_store

import (
	"lightning/app/plain/log"
)

var store map[string]string

func InitClient() error {
	if store == nil {
		store = make(map[string]string)
	}

	return nil
}

func Set[K ~string, V ~string](k K, v V) error {
	log.Debug("[dist_store] setting", k, "to", v)
	store[string(k)] = string(v)
	return nil
}

func Get[K ~string](k K) (string, bool) {
	_, ok := store[string(k)]
	return store[string(k)], ok
}

func Close() {}
