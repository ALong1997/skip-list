package skiplist

import (
	"golang.org/x/exp/constraints"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type (
	// ConcurrentSkipList is thread-safe.
	ConcurrentSkipList[O constraints.Ordered, T any] struct {
		level, maxLevel, cap atomic.Uint32
		head                 *concurrentNode[O, T]
		r                    *rand.Rand
		nodeCache            sync.Pool

		deleteRWMutex sync.RWMutex
		putMap        sync.Map
	}

	concurrentNode[O constraints.Ordered, T any] struct {
		node[O, T]

		sync.RWMutex
	}
)

func NewConcurrentSkipList[O constraints.Ordered, T any](maxLevel uint32) *ConcurrentSkipList[O, T] {
	if maxLevel == 0 {
		return nil
	}

	var csl = &ConcurrentSkipList[O, T]{
		head:      &concurrentNode[O, T]{node: node[O, T]{nextNodes: make([]*node[O, T], 1)}},
		r:         rand.New(rand.NewSource(time.Now().Unix())),
		nodeCache: sync.Pool{New: func() any { return &node[O, T]{} }},
	}

	csl.level.Store(1)
	csl.maxLevel.Store(maxLevel)

	return csl
}
