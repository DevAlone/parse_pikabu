package helpers

import (
	"sync"
	"time"
)

// IDLocker is used to lock items with the same id
type IDLocker struct {
	locks map[uint64]bool
	mutex sync.Mutex
}

// Lock locks item with id
func (l *IDLocker) Lock(id uint64) {
	found := true
	for found {
		l.mutex.Lock()
		_, found = l.locks[id]
		if !found {
			l.locks[id] = true
			l.mutex.Unlock()
			return
		}
		l.mutex.Unlock()

		time.Sleep(10 * time.Millisecond)
	}
}

// Unlock unlocks item with id
func (l *IDLocker) Unlock(id uint64) {
	l.mutex.Lock()
	delete(l.locks, id)
	l.mutex.Unlock()
}
