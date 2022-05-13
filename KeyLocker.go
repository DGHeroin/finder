package finder

import "sync"

type (
    KeyLocker struct {
        sync.Mutex
        count int
    }
    Locker struct {
        lock  sync.Mutex
        locks map[string]*KeyLocker
    }
)

func NewLocker() *Locker {
    return &Locker{
        lock:  sync.Mutex{},
        locks: make(map[string]*KeyLocker),
    }
}

func (l *Locker) Lock(key string) {
    l.lock.Lock()

    keyLocker := l.locks[key]
    if keyLocker == nil {
        keyLocker = &KeyLocker{}
        l.locks[key] = keyLocker
    }
    keyLocker.count++
    l.lock.Unlock()

    keyLocker.Lock()
}
func (l *Locker) UnLock(key string) {
    l.lock.Lock()

    keyLocker := l.locks[key]
    keyLocker.Unlock()
    keyLocker.count--
    if keyLocker.count == 0 {
        delete(l.locks, key)
    }

    l.lock.Unlock()
}

