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

		deleteRWMutex sync.RWMutex // Only locked by ConcurrentSkipList.Delete
		keyPutMap     sync.Map     // The same key is writen by ConcurrentSkipList.Put
	}

	concurrentNode[O constraints.Ordered, T any] struct {
		*KvPair[O, T]
		nextNodes []*concurrentNode[O, T]

		sync.RWMutex
	}
)

func NewConcurrentSkipList[O constraints.Ordered, T any](maxLevel uint32) *ConcurrentSkipList[O, T] {
	if maxLevel == 0 {
		return nil
	}

	var csl = &ConcurrentSkipList[O, T]{
		head:      &concurrentNode[O, T]{nextNodes: make([]*concurrentNode[O, T], 1)},
		r:         rand.New(rand.NewSource(time.Now().Unix())),
		nodeCache: sync.Pool{New: func() any { return &concurrentNode[O, T]{} }},
	}

	csl.level.Store(1)
	csl.maxLevel.Store(maxLevel)

	return csl
}

func (csl *ConcurrentSkipList[O, T]) Level() uint32 {
	csl.deleteRWMutex.RLock()
	defer csl.deleteRWMutex.RUnlock()
	return csl.level.Load()
}

func (csl *ConcurrentSkipList[O, T]) Cap() uint32 {
	csl.deleteRWMutex.RLock()
	defer csl.deleteRWMutex.RUnlock()
	return csl.cap.Load()
}

func (csl *ConcurrentSkipList[O, T]) Get(key O) (val T, exist bool) {
	if csl.Level() == 0 {
		return
	}

	csl.deleteRWMutex.RLock()
	defer csl.deleteRWMutex.RUnlock()

	if n := csl.get(key); n != nil {
		return n.val, true
	}
	return
}

func (csl *ConcurrentSkipList[O, T]) Put(key O, val T) {
	if csl.Level() == 0 {
		return
	}

	csl.deleteRWMutex.RLock()
	defer csl.deleteRWMutex.RUnlock()

	keyPutMutex := csl.getKeyPutMutex(key)
	keyPutMutex.Lock()
	defer keyPutMutex.Unlock()

	n := csl.get(key)
	if n != nil {
		// update
		n.val = val
		return
	}

	// randomly determined level
	var randL = csl.randLevel()

	// grow
	csl.grow(randL + 1)

	// new node
	n, _ = csl.nodeCache.Get().(*concurrentNode[O, T])
	n.KvPair = newKvPair(key, val)
	n.nextNodes = make([]*concurrentNode[O, T], randL+1)

	n.Lock()
	defer n.Unlock()

	var preNode *concurrentNode[O, T]
	current := csl.head
	for l := len(current.nextNodes) - 1; l >= 0; l-- {
		for current.nextNodes[l] != nil && current.key < key {
			// search to the right
			current = current.nextNodes[l]
		}

		if current != preNode {
			// every level preNode needs lock, but don't lock a same node
			current.RLock()
			defer current.RUnlock()
			preNode = current
		}

		// insert
		n.nextNodes[l] = current.nextNodes[l]
		current.nextNodes[l] = n

		// search down
	}

	csl.cap.Add(1)
}

// Delete is mutually exclusive with others
func (csl *ConcurrentSkipList[O, T]) Delete(key O) {
	if csl.Level() == 0 {
		return
	}

	csl.deleteRWMutex.Lock()
	defer csl.deleteRWMutex.Unlock()

	var deleteNode *concurrentNode[O, T]
	current := csl.head
	for l := len(current.nextNodes) - 1; l >= 0; l-- {
		for current.nextNodes[l] != nil && current.nextNodes[l].key < key {
			// search to the right
			current = current.nextNodes[l]
		}

		if current.nextNodes[l] != nil && current.nextNodes[l].key == key {
			// delete
			if deleteNode == nil {
				deleteNode = current.nextNodes[l]
			}
			current.nextNodes[l] = current.nextNodes[l].nextNodes[l]
		}

		// search down
	}

	if deleteNode == nil {
		// not exist
		return
	}
	deleteNode.nextNodes = nil
	csl.nodeCache.Put(deleteNode)

	// cut
	csl.cut()

	csl.cap.CompareAndSwap(csl.Cap(), csl.Cap()-1)
}

/*
// Range searches the *KvPair of key in [start, end].

	func (csl *ConcurrentSkipList[O, T]) Range(start, end O) []*KvPair[O, T] {
		if csl.Level() == 0 {
			return nil
		}

		var res = make([]*KvPair[O, T], 0)

		// starting point
		ceilingNode := csl.ceil(start)
		if ceilingNode == nil {
			return res
		}

		// range
		for n := ceilingNode; n != nil && n.key <= end; n = n.nextNodes[0] {
			res = append(res, newKvPair(n.key, n.val))
		}
		return res
	}

// Ceil returns *KvPair of the least key greater than or equal to target.

	func (csl *ConcurrentSkipList[O, T]) Ceil(target O) (*KvPair[O, T], bool) {
		if csl.Level() == 0 {
			return nil, false
		}

		if ceilingNode := csl.ceil(target); ceilingNode != nil {
			return newKvPair(ceilingNode.key, ceilingNode.val), true
		}
		return nil, false
	}

// Floor returns *KvPair of the greatest key less than or equal to target.

	func (csl *ConcurrentSkipList[O, T]) Floor(target O) (*KvPair[O, T], bool) {
		if csl.Level() == 0 {
			return nil, false
		}

		if floorNode := csl.floor(target); floorNode != csl.head.nextNodes[0] {
			return newKvPair(floorNode.key, floorNode.val), true
		}
		return nil, false
	}
*/
func (csl *ConcurrentSkipList[O, T]) get(key O) *concurrentNode[O, T] {
	if csl.Level() == 0 {
		return nil
	}

	var preNode *concurrentNode[O, T]
	current := csl.head
	for l := len(current.nextNodes) - 1; l >= 0; l-- {
		for current.nextNodes[l] != nil && current.nextNodes[l].key < key {
			// search to the right
			current = current.nextNodes[l]
		}

		if current != preNode {
			// every level preNode needs lock, but don't lock a same node
			current.RLock()
			defer current.RUnlock()
			preNode = current
		}

		if current.nextNodes[l] != nil && current.nextNodes[l].key == key {
			// exist
			return current.nextNodes[l]
		}

		// search down
	}
	// not exist
	return nil
}

func (csl *ConcurrentSkipList[O, T]) ceil(target O) *concurrentNode[O, T] {
	if csl.Level() == 0 {
		return nil
	}

	current := csl.head
	for l := len(current.nextNodes) - 1; l >= 0; l-- {
		for current.nextNodes[l] != nil && current.nextNodes[l].key < target {
			// search to the right
			current = current.nextNodes[l]
		}

		if current.nextNodes[l] != nil && current.nextNodes[l].key == target {
			// equal
			return current.nextNodes[l]
		}

		// search down
	}
	// current.nextNodes[0] is ceil || current.nextNodes[0] == nil(tail node means ceil is not exist)
	return current.nextNodes[0]
}

func (csl *ConcurrentSkipList[O, T]) floor(target O) *concurrentNode[O, T] {
	if csl.Level() == 0 {
		return nil
	}

	current := csl.head
	for l := len(current.nextNodes) - 1; l >= 0; l-- {
		for current.nextNodes[l] != nil && current.nextNodes[l].key < target {
			// search to the right
			current = current.nextNodes[l]
		}

		if current.nextNodes[l] != nil && current.nextNodes[l].key == target {
			// equal
			return current.nextNodes[l]
		}

		// search down
	}
	// current is floor || current == csl.current.nextNodes[0](head node means floor is not exist)
	return current
}

func (csl *ConcurrentSkipList[O, T]) randLevel() uint32 {
	var randL uint32
	for rand.Intn(2) == 0 && randL < csl.maxLevel.Load() {
		randL++
	}
	return randL
}

func (csl *ConcurrentSkipList[O, T]) grow(newL uint32) {
	if csl.Level() < newL {
		csl.head.Lock()
		csl.head.nextNodes = append(csl.head.nextNodes, make([]*concurrentNode[O, T], newL-csl.Level())...)
		csl.level.Store(newL)
		csl.head.Unlock()
	}
}

// ConcurrentSkipList.cut only be called by ConcurrentSkipList.Delete, so ConcurrentSkipList.cut needn't lock
func (csl *ConcurrentSkipList[O, T]) cut() {
	var dif uint32
	for l := csl.Level() - 1; l > 0; l-- {
		if csl.head.nextNodes[l] != nil {
			break
		}
		dif++
	}
	csl.head.nextNodes = csl.head.nextNodes[:csl.Level()-dif]

	csl.level.Add(^dif + 1)
}

func (csl *ConcurrentSkipList[O, T]) getKeyPutMutex(key O) *sync.Mutex {
	rawMutex, _ := csl.keyPutMap.LoadOrStore(key, &sync.Mutex{})
	mutex, _ := rawMutex.(*sync.Mutex)
	return mutex
}
