package sdb

import "sync"

// database alias cacher.
type _sdbCache struct {
	mux     sync.RWMutex
	cache   map[string]*Sdb
	nameMap map[string]string
}

const defaultDbName = "default"

// add database alias with original Name.
func (ac *_sdbCache) add(sdb *Sdb) (added bool) {
	ac.mux.Lock()
	defer ac.mux.Unlock()
	// 不能添加同文件的数据库
	if _, ok := ac.nameMap[sdb.GetPath()]; ok {
		return false
	}
	// 不能添加同名的数据库
	if _, ok := ac.cache[sdb.Name]; !ok {
		ac.cache[sdb.Name] = sdb
		ac.nameMap[sdb.GetPath()] = sdb.Name
		added = true
	}
	return
}

// get database alias if cached.
func (ac *_sdbCache) get(name string) (sdb *Sdb, ok bool) {
	ac.mux.RLock()
	defer ac.mux.RUnlock()
	sdb, ok = ac.cache[name]
	return
}

// get default alias.
func (ac *_sdbCache) getDefault() (*Sdb, bool) {
	return ac.get(defaultDbName)
}
