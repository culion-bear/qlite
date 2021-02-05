package localTime

import "sync"

type Timer struct {
	timer	int64
	lock	sync.RWMutex
}