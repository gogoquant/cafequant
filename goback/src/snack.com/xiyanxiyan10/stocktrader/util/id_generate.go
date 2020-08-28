package util

import (
	"fmt"
	"sync"
	"time"
)

// IDGen ...
type IDGen struct {
	Prefix string
	ID     int64
	*sync.Mutex
}

// NewIDGen ...
func NewIDGen(prefix string) *IDGen {
	date := time.Now().Format("2006-01-02")
	return &IDGen{
		Mutex:  new(sync.Mutex),
		Prefix: fmt.Sprintf("%s.%s", prefix, date),
		ID:     0,
	}
}

// Get ...
func (gen *IDGen) Get() string {
	gen.Lock()
	defer gen.Unlock()
	gen.ID++
	return fmt.Sprintf("%s.%d", gen.Prefix, gen.ID)
}
