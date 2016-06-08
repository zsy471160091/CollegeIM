package process

import (
	"sync"
)

type clientList struct {
	m       map[string](chan<- []byte)
	rwmutex sync.RWMutex
}

func newClientList(clientNum int) *clientList {
	return &clientList{
		m: make(map[string](chan<- []byte), clientNum),
	}
}

func (cl *clientList) get(key string) chan<- []byte {
	cl.rwmutex.RLock()
	defer cl.rwmutex.RUnlock()
	return cl.m[key]
}

func (cl *clientList) add(key string, elem chan<- []byte) (chan<- []byte, bool) {
	cl.rwmutex.Lock()
	defer cl.rwmutex.Unlock()
	oldElem := cl.m[key]
	cl.m[key] = elem
	return oldElem, true
}

func (cl *clientList) remove(key string) chan<- []byte {
	cl.rwmutex.Lock()
	defer cl.rwmutex.Unlock()
	oldElem := cl.m[key]
	delete(cl.m, key)
	return oldElem
}

func (cl *clientList) clear() {
	cl.rwmutex.Lock()
	defer cl.rwmutex.Unlock()
	cl.m = make(map[string](chan<- []byte))
}

func (cl *clientList) len() int {
	cl.rwmutex.RLock()
	defer cl.rwmutex.RUnlock()
	return len(cl.m)
}

func (cl *clientList) contains(key string) bool {
	cl.rwmutex.RLock()
	defer cl.rwmutex.RUnlock()
	_, ok := cl.m[key]
	return ok
}
