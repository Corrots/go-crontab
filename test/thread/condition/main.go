package main

import (
	"sync"
)

var (
	cond    sync.Cond
	rwmutex sync.RWMutex
	mutex   sync.Mutex
)

func main() {
	cond.Wait()
	cond.Signal()
	rwmutex.RLock()
	rwmutex.RUnlock()
	mutex.Lock()
	mutex.Unlock()
}
