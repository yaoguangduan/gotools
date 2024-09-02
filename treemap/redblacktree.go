package algo

import (
	"cmp"
	"iter"
)

type CompareFunc[T any] func(a, b T) int

type color uint8

const (
	black color = iota
	red
)
const treeFreeMaxSize = 16

type treeNode[K, V any] struct {
	key    K
	value  V
	left   *treeNode[K, V]
	right  *treeNode[K, V]
	parent *treeNode[K, V]
	color  color
}

type RedBlackTree[K, V any] struct {
	root     *treeNode[K, V]
	cmp      CompareFunc[K]
	length   int
	freeSize int
	frees    *treeNode[K, V]
}

func New[K cmp.Ordered, V any]() *RedBlackTree[K, V] {
	return NewWithCmpFunc[K, V](func(a, b K) int {
		return cmp.Compare(a, b)
	})
}
func NewWithCmpFunc[K, V any](cmpFn CompareFunc[K]) *RedBlackTree[K, V] {
	return &RedBlackTree[K, V]{cmp: cmpFn}
}
func (t *RedBlackTree[K, V]) Len() int {
	return t.length
}
func (t *RedBlackTree[K, V]) Get(key K) V {
	v, _ := t.TryGet(key)
	return v
}
func (t *RedBlackTree[K, V]) TryGet(key K) (V, bool) {
	p := t.root
	for p != nil {
		cp := t.cmp(key, p.key)
		if cp < 0 {
			p = p.left
		} else if cp > 0 {
			p = p.right
		} else {
			return p.value, true
		}
	}
	return *new(V), false
}
func (t *RedBlackTree[K, V]) Put(key K, value V) (V, bool) {
	if t.root == nil {
		t.root = t.alloc(key, value)
		t.length = 1
		return *new(V), false
	}
	var cp = 0
	var p = t.root
	var parent *treeNode[K, V]
	for p != nil {
		parent = p
		cp = t.cmp(key, p.key)
		if cp < 0 {
			p = p.left
		} else if cp > 0 {
			p = p.right
		} else {
			old := p.value
			p.value = value
			return old, true
		}
	}
	node := t.alloc(key, value)
	if cp < 0 {
		parent.left = node
		parent.left.parent = parent
	} else {
		parent.right = node
		parent.right.parent = parent
	}
	t.fixWhenInsert(node)
	t.length++
	return *new(V), false
}
func (t *RedBlackTree[K, V]) Delete(key K) (V, bool) {
	node := t.getNode(key)
	if node == nil {
		return *new(V), false
	}
	v := node.value
	t.deleteNode(node)
	return v, true
}
func (t *RedBlackTree[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		t.visit(t.root, yield, alwaysTrueCmp[K]{})
	}
}

func (t *RedBlackTree[K, V]) Keys() []K {
	ks := make([]K, 0)
	for k, _ := range t.Iter() {
		ks = append(ks, k)
	}
	return ks
}

type treeInRangeCmp[T any] interface {
	inRange(target T) bool
}
type alwaysTrueCmp[T any] struct {
}

func (at alwaysTrueCmp[T]) inRange(T) bool {
	return true
}

type doubleDirectionCmp[T any] struct {
	beg T
	end T
	cmp func(T, T) int
}

func (vc doubleDirectionCmp[T]) inRange(target T) bool {
	return vc.cmp(vc.beg, target) <= 0 && vc.cmp(vc.end, target) > 0
}

type singleDirectionCmp[T any] struct {
	val T
	beg bool
	cmp func(T, T) int
}

func (vc singleDirectionCmp[T]) inRange(target T) bool {
	if vc.beg {
		return vc.cmp(vc.val, target) <= 0
	} else {
		return vc.cmp(vc.val, target) >= 0
	}
}

func (t *RedBlackTree[K, V]) IterRange(beg, end K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		t.visit(t.root, yield, doubleDirectionCmp[K]{beg: beg, end: end, cmp: t.cmp})
	}
}

func (t *RedBlackTree[K, V]) IterBE(beg K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		t.visit(t.root, yield, singleDirectionCmp[K]{val: beg, beg: true, cmp: t.cmp})
	}
}

func (t *RedBlackTree[K, V]) IterLE(end K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		t.visit(t.root, yield, singleDirectionCmp[K]{val: end, cmp: t.cmp})
	}
}

func (t *RedBlackTree[K, V]) fixWhenInsert(node *treeNode[K, V]) {
	node.color = red
	for node != nil && node != t.root && node.parent.color == red {
		if parentOf(node) == leftOf(parentOf(parentOf(node))) {
			r := rightOf(parentOf(parentOf(node)))
			if colorOf(r) == red {
				setColor(parentOf(node), black)
				setColor(r, black)
				setColor(parentOf(parentOf(node)), red)
				node = parentOf(parentOf(node))
			} else {
				if node == rightOf(parentOf(node)) {
					node = parentOf(node)
					t.rotateLeft(node)
				}
				setColor(parentOf(node), black)
				setColor(parentOf(parentOf(node)), red)
				t.rotateRight(parentOf(parentOf(node)))
			}
		} else {
			u := leftOf(parentOf(parentOf(node)))
			if colorOf(u) == red {
				setColor(parentOf(parentOf(node)), red)
				setColor(u, black)
				setColor(parentOf(node), black)
				node = parentOf(parentOf(node))
			} else {
				if node == leftOf(parentOf(node)) {
					node = parentOf(node)
					t.rotateRight(node)
				}
				setColor(parentOf(parentOf(node)), red)
				setColor(parentOf(node), black)
				t.rotateLeft(parentOf(parentOf(node)))
			}
		}
	}
}

func leftOf[K, V any](node *treeNode[K, V]) *treeNode[K, V] {
	if node == nil {
		return nil
	}
	return node.left
}
func rightOf[K, V any](node *treeNode[K, V]) *treeNode[K, V] {
	if node == nil {
		return nil
	}
	return node.right
}

func parentOf[K, V any](node *treeNode[K, V]) *treeNode[K, V] {
	if node == nil {
		return nil
	}
	return node.parent
}

func colorOf[K, V any](node *treeNode[K, V]) color {
	if node == nil {
		return black
	}
	return node.color
}
func setColor[K, V any](node *treeNode[K, V], color color) {
	if node == nil {
		return
	}
	node.color = color
}
func (t *RedBlackTree[K, V]) rotateLeft(node *treeNode[K, V]) {
	if node == nil {
		return
	}
	right := rightOf(node)
	if right == nil {
		return
	}
	node.right = leftOf(right)
	if right.left != nil {
		right.left.parent = node
	}
	p := parentOf(node)
	right.parent = p
	if p == nil {
		t.root = right
	} else if node == leftOf(p) {
		p.left = right
	} else {
		p.right = right
	}
	right.left = node
	node.parent = right
}
func (t *RedBlackTree[K, V]) rotateRight(node *treeNode[K, V]) {
	if node == nil {
		return
	}
	left := leftOf(node)
	if left == nil {
		return
	}
	node.left = left.right
	if left.right != nil {
		left.right.parent = node
	}
	p := parentOf(node)
	left.parent = p
	if p == nil {
		t.root = left
	} else if node == leftOf(p) {
		p.left = left
	} else {
		p.right = left
	}
	left.right = node
	node.parent = left
}

func (t *RedBlackTree[K, V]) getNode(key K) *treeNode[K, V] {
	var tmp = t.root
	for tmp != nil {
		cp := t.cmp(key, tmp.key)
		if cp < 0 {
			tmp = tmp.left
		} else if cp > 0 {
			tmp = tmp.right
		} else {
			return tmp
		}
	}
	return nil
}

func (t *RedBlackTree[K, V]) deleteNode(node *treeNode[K, V]) {
	t.length--
	if node.left != nil && node.right != nil {
		s := t.successor(node)
		node.key = s.key
		node.value = s.value
		node = s
	}
	var replaceNode = node.left
	if replaceNode == nil {
		replaceNode = node.right
	}
	if replaceNode != nil {
		replaceNode.parent = node.parent
		if node.parent != nil {
			if node == node.parent.left {
				node.parent.left = replaceNode
			} else if node == node.parent.right {
				node.parent.right = replaceNode
			}
		} else {
			t.root = replaceNode
		}
		node.parent, node.left, node.right = nil, nil, nil
		t.free(node)
		if node.color == black {
			t.fixWhenDelete(replaceNode)
		}
	} else if node.parent == nil {
		t.root = nil
	} else {
		if node.color == black {
			t.fixWhenDelete(node)
		}
		if node.parent != nil {
			if node == node.parent.left {
				node.parent.left = nil
			} else if node.parent.right == node {
				node.parent.right = nil
			}
			node.parent = nil
			t.free(node)
		}
	}

}

func (t *RedBlackTree[K, V]) successor(node *treeNode[K, V]) *treeNode[K, V] {
	if node == nil {
		return nil
	} else if node.right != nil {
		r := node.right
		for r.left != nil {
			r = r.left
		}
		return r
	} else {
		p := node.parent
		ch := node
		for p != nil && ch == p.right {
			ch = p
			p = p.parent
		}
		return p
	}

}

func (t *RedBlackTree[K, V]) fixWhenDelete(node *treeNode[K, V]) {
	for node != t.root && node.color == black {
		if node == leftOf(parentOf(node)) {
			uncle := rightOf(parentOf(node))
			if colorOf(uncle) == red {
				setColor(uncle, black)
				setColor(parentOf(node), red)
				t.rotateLeft(parentOf(node))
				uncle = rightOf(parentOf(node))
			}
			if colorOf(leftOf(uncle)) == black && colorOf(rightOf(uncle)) == black {
				setColor(uncle, red)
				node = parentOf(node)
			} else {
				if colorOf(rightOf(uncle)) == black {
					setColor(leftOf(uncle), black)
					setColor(uncle, red)
					t.rotateRight(uncle)
					uncle = rightOf(parentOf(node))
				}
				setColor(uncle, colorOf(parentOf(node)))
				setColor(parentOf(node), black)
				setColor(rightOf(uncle), black)
				t.rotateLeft(parentOf(node))
				node = t.root
			}
		} else {
			uncle := leftOf(parentOf(node))
			if colorOf(uncle) == red {
				setColor(uncle, black)
				setColor(parentOf(node), red)
				t.rotateRight(parentOf(node))
				uncle = leftOf(parentOf(node))
			}

			if colorOf(leftOf(uncle)) == black && colorOf(rightOf(uncle)) == black {
				setColor(uncle, red)
				node = parentOf(node)
			} else {
				if colorOf(leftOf(uncle)) == black {
					setColor(rightOf(uncle), black)
					setColor(uncle, red)
					t.rotateLeft(uncle)
					uncle = leftOf(parentOf(node))
				}
				setColor(uncle, colorOf(parentOf(node)))
				setColor(parentOf(node), black)
				setColor(leftOf(uncle), black)
				t.rotateRight(parentOf(node))
				node = t.root
			}
		}
		setColor(node, black)
	}
}

func (t *RedBlackTree[K, V]) visit(root *treeNode[K, V], yield func(K, V) bool, ir treeInRangeCmp[K]) bool {
	if root == nil {
		return true
	}
	if !t.visit(root.left, yield, ir) {
		return false
	}

	if ir.inRange(root.key) {
		if !yield(root.key, root.value) {
			return false
		}
	}
	if !t.visit(root.right, yield, ir) {
		return false
	}

	return true
}

func (t *RedBlackTree[K, V]) alloc(key K, value V) *treeNode[K, V] {
	if t.frees == nil {
		return &treeNode[K, V]{key: key, value: value}
	}
	ret := t.frees
	t.frees = t.frees.parent
	ret.parent = nil
	ret.key = key
	ret.value = value
	t.freeSize--
	return ret
}
func (t *RedBlackTree[K, V]) free(n *treeNode[K, V]) {
	if t.freeSize >= treeFreeMaxSize {
		return
	}
	t.freeSize++
	n.parent = t.frees
	t.frees = n
}
