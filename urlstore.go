/*

Line-based url storage for now

NB uint32s should become uint64 at some point

*/

package grepurl

type URLStore interface {
	addURL(url string, trigrams []string) uint32
	getURL(id uint32) string
	getCardinality() uint32
}

type MemoryURLStore struct {
	urls []string
}

func NewMemoryURLStore() *MemoryURLStore {
	return &MemoryURLStore{urls: make([]string, 0)}
}

func (us *MemoryURLStore) getCardinality() uint32 {
	return uint32(len(us.urls))
}

func (us *MemoryURLStore) addURL(url string, trigrams []string) uint32 {
	us.urls = append(us.urls, url)
	return us.getCardinality() - 1
}

func (us *MemoryURLStore) getURL(id uint32) string {
	return us.urls[id]
}
