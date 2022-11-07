// 根据key单独加锁
package kmutex

import (
	"sync"
)

type KMutex struct {
	c *sync.Cond
	l sync.Locker
	s map[any]struct{}
}

// Create new KMutex
func New() *KMutex {
	l := sync.Mutex{}
	return &KMutex{
		c: sync.NewCond(&l),
		l: &l,
		s: make(map[any]struct{}),
	}
}

func (km *KMutex) locked(key any) (ok bool) {
	_, ok = km.s[key]
	return
}

// Unlock KMutex by unique ID
// 每次Unlock都会唤醒所有Wait
func (km *KMutex) Unlock(key any) {
	km.l.Lock()
	defer km.l.Unlock()
	delete(km.s, key)
	km.c.Broadcast()
}

// Lock KMutex by unique ID
func (km *KMutex) Lock(key any) {
	km.l.Lock()
	defer km.l.Unlock()
	for km.locked(key) {
		km.c.Wait()
	}
	km.s[key] = struct{}{}
}
