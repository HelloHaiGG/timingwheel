package utils

import (
	"sync"
)

var (
	PbPool *PBMessagePool
)

func init() {
	PbPool = NewPbMessagePool()
}

type PBMessagePool struct {
	pool sync.Pool
}

func NewPbMessagePool() *PBMessagePool {
	return &PBMessagePool{pool: sync.Pool{New: func() interface{} {
		return nil
	}}}
}

func (p *PBMessagePool) Get() interface{} {
	return nil
}
