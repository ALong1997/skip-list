package skiplist

import (
	"math/rand"
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

type (
	// SkipList is not thread-safe.
	SkipList[O constraints.Ordered, T any] struct {
		level, maxLevel, cap uint32
		head                 *node[O, T]
		r                    *rand.Rand
		nodeCache            sync.Pool
	}

	node[O constraints.Ordered, T any] struct {
		*KvPair[O, T]
		nextNodes []*node[O, T]
	}
)

func NewSkipList[O constraints.Ordered, T any](maxLevel uint32) *SkipList[O, T] {
	if maxLevel <= 0 {
		return nil
	}

	return &SkipList[O, T]{
		level:     1,
		maxLevel:  maxLevel,
		cap:       0,
		head:      &node[O, T]{nextNodes: make([]*node[O, T], 1)},
		r:         rand.New(rand.NewSource(time.Now().Unix())),
		nodeCache: sync.Pool{New: func() any { return &node[O, T]{} }},
	}
}

func (sl *SkipList[O, T]) Level() uint32 {
	return sl.level
}

func (sl *SkipList[O, T]) Cap() uint32 {
	return sl.cap
}

func (sl *SkipList[O, T]) Get(key O) (val T, exist bool) {
	if sl.Level() == 0 {
		return
	}

	if n := sl.get(key); n != nil {
		return n.val, true
	}
	return
}

func (sl *SkipList[O, T]) Put(key O, val T) {
	if sl.Level() == 0 {
		return
	}

	n := sl.get(key)
	if n != nil {
		// update
		n.val = val
		return
	}

	// randomly determined level
	var randL = sl.randLevel()

	// grow
	if sl.Level() < randL+1 {
		sl.head.nextNodes = append(sl.head.nextNodes, make([]*node[O, T], randL+1-sl.Level())...)
		sl.level = randL + 1
	}

	// new node
	n, _ = sl.nodeCache.Get().(*node[O, T])
	n.KvPair = newKvPair(key, val)
	n.nextNodes = make([]*node[O, T], randL+1)

	current := sl.head
	for l := len(current.nextNodes) - 1; l >= 0; l-- {
		for current.nextNodes[l] != nil && current.key < key {
			// search to the right
			current = current.nextNodes[l]
		}

		// insert
		n.nextNodes[l] = current.nextNodes[l]
		current.nextNodes[l] = n

		// search down
	}

	sl.cap++
}

func (sl *SkipList[O, T]) Del(key O) {
	if sl.Level() == 0 {
		return
	}

	var deleteNode *node[O, T]
	current := sl.head
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
	sl.nodeCache.Put(deleteNode)

	// cut
	var dif uint32
	for l := sl.Level() - 1; l >= 0; l-- {
		if sl.head.nextNodes[l] != nil {
			break
		}
		dif++
	}
	sl.head.nextNodes = sl.head.nextNodes[:sl.Level()-dif]

	sl.level -= dif
	sl.cap--
}

// Range searches the *KvPair of key in [start, end].
func (sl *SkipList[O, T]) Range(start, end O) []*KvPair[O, T] {
	if sl.Level() == 0 {
		return nil
	}

	var res = make([]*KvPair[O, T], 0)

	// starting point
	ceilingNode := sl.ceil(start)
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
func (sl *SkipList[O, T]) Ceil(target O) (*KvPair[O, T], bool) {
	if sl.Level() == 0 {
		return nil, false
	}

	if ceilingNode := sl.ceil(target); ceilingNode != nil {
		return newKvPair(ceilingNode.key, ceilingNode.val), true
	}
	return nil, false
}

// Floor returns *KvPair of the greatest key less than or equal to target.
func (sl *SkipList[O, T]) Floor(target O) (*KvPair[O, T], bool) {
	if sl.Level() == 0 {
		return nil, false
	}

	if floorNode := sl.floor(target); floorNode != sl.head.nextNodes[0] {
		return newKvPair(floorNode.key, floorNode.val), true
	}
	return nil, false
}

func (sl *SkipList[O, T]) get(key O) *node[O, T] {
	if sl.Level() == 0 {
		return nil
	}

	current := sl.head
	for l := len(current.nextNodes) - 1; l >= 0; l-- {
		for current.nextNodes[l] != nil && current.nextNodes[l].key < key {
			// search to the right
			current = current.nextNodes[l]
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

func (sl *SkipList[O, T]) ceil(target O) *node[O, T] {
	if sl.Level() == 0 {
		return nil
	}

	current := sl.head
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

func (sl *SkipList[O, T]) floor(target O) *node[O, T] {
	if sl.Level() == 0 {
		return nil
	}

	current := sl.head
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
	// current is floor || current == sl.current.nextNodes[0](head node means floor is not exist)
	return current
}

func (sl *SkipList[O, T]) randLevel() uint32 {
	var randL uint32
	for rand.Intn(2) == 0 && randL < sl.maxLevel {
		randL++
	}
	return randL
}
