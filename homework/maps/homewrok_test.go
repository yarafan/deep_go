package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go
type Comparator[K any] func(a, b K) int

type Node[K comparable, V any] struct {
	parent *Node[K, V]
	left   *Node[K, V]
	right  *Node[K, V]
	key    K
	value  V
}

type OrderedMap[K comparable, V any] struct {
	size       int
	root       *Node[K, V]
	comparator Comparator[K]
}

func NewOrderedMap[K comparable, V any](cmp Comparator[K]) OrderedMap[K, V] {
	return OrderedMap[K, V]{comparator: cmp}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	if m.root == nil {
		m.root = &Node[K, V]{key: key, value: value}

		m.size++

		return
	}

	root := m.root
	parent := (*Node[K, V])(nil)

	for root != nil {
		cmp := m.comparator(key, root.key)

		if cmp == 0 {
			root.value = value

			return
		}

		parent = root
		if cmp < 0 {
			root = root.left
		} else {
			root = root.right
		}
	}

	newNode := &Node[K, V]{key: key, value: value, parent: parent}
	if m.comparator(key, parent.key) < 0 {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	m.size++
}

func (m *OrderedMap[K, V]) Erase(key K) {
	node := m.findNode(key)
	if node == nil {
		return
	}

	if node.right == nil {
		m.transplant(node, node.left)
	} else {
		successor := node.right

		for successor.left != nil {
			successor = successor.left
		}

		if successor.parent != node {
			m.transplant(successor, successor.right)
			successor.right = node.right
			if successor.right != nil {
				successor.right.parent = successor
			}
		}

		m.transplant(node, successor)
		successor.left = node.left
		if successor.left != nil {
			successor.left.parent = successor
		}

	}

	m.size--
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	node := m.findNode(key)

	return node != nil
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	stack := make([]*Node[K, V], 0)
	curr := m.root

	for curr != nil || len(stack) > 0 {
		for curr != nil {
			stack = append(stack, curr)
			curr = curr.left
		}

		curr = stack[len(stack)-1]

		action(curr.key, curr.value)

		curr = curr.right
		stack = stack[:len(stack)-1]
	}
}

func (m *OrderedMap[K, V]) findNode(key K) *Node[K, V] {
	root := m.root

	for root != nil {
		cmp := m.comparator(key, root.key)
		if cmp == 0 {
			return root
		}

		if cmp < 0 {
			root = root.left
		} else {
			root = root.right
		}
	}

	return nil
}

func (m *OrderedMap[K, V]) transplant(predecessor *Node[K, V], successor *Node[K, V]) {
	if predecessor.parent == nil {
		m.root = successor

		return
	}

	if predecessor == predecessor.parent.left {
		predecessor.parent.left = successor
	} else {
		predecessor.parent.right = successor
	}

	if successor != nil {
		successor.parent = predecessor.parent
	}
}

func intComparator(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int](intComparator)
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
