// Stand-in for some kind of key-value distributed storage, like etcd
package dist_store

var store map[string]string

func Set[K ~string, V ~string](k K, v V) error {
	store[string(k)] = string(v)
	return nil
}

func Get[K ~string](k K) (string, bool) {
	_, ok := store[string(k)]
	return store[string(k)], ok
}
