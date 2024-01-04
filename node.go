package skiplist

import "golang.org/x/exp/constraints"

type (
	node[O constraints.Ordered, T any] struct {
		*KvPair[O, T]
		nextNodes []*node[O, T]
	}

	KvPair[O constraints.Ordered, T any] struct {
		key O
		val T
	}
)

func (kv *KvPair[O, T]) Key() (key O) {
	return kv.key
}

func (kv *KvPair[O, T]) Val() (val T) {
	return kv.val
}

func newKvPair[O constraints.Ordered, T any](key O, val T) *KvPair[O, T] {
	return &KvPair[O, T]{
		key: key,
		val: val,
	}
}

func newNode[O constraints.Ordered, T any](key O, val T, nextNodes []*node[O, T]) *node[O, T] {
	return &node[O, T]{
		KvPair:    newKvPair(key, val),
		nextNodes: nextNodes,
	}
}
