// Package heap provide the main Data container
// heap for the log analyse tool
// every node of heap contain value and count
// heap also maintain a map of value and node
// to make a quick search for duplicate value
package heap

// HeapType is used to specify a heap is min-heap or max-heap
type HeapType int32

const (
	// MaxHeap represents max-heap
	MaxHeap HeapType = 1
	// MinHeap represents min-heap
	MinHeap HeapType = 2
)

// Node is the base type of each heap node
// it contain the value for the node and total
// count for the same value
type Node struct {
	Val   uint64
	Count uint64
}

// Heap is the main struct for the heap container
// it contains the node list which stores the real Data
// Type is used to specify whether the heap is min-heap or max-heap
// NodeMap is used to make a quick index for duplicate values
// set the root index of heap to 1 to find the sub node index easily
type Heap struct {
	Data       []*Node
	TotalCount uint64
	Type       HeapType
	NodeMap    map[uint64]*Node
}

// NewHeap create a new heap container
// it leave the first item empty to make index calculate simpler
func NewHeap(capacity int, heapType HeapType) *Heap {
	return &Heap{
		make([]*Node, 1, capacity),
		0,
		heapType,
		make(map[uint64]*Node)}
}

// Insert is to insert a value to a heap
// it first check if the val is a duplicate value
// if so, just add the node count
// if not, just add a new node
func (h *Heap) Insert(val uint64) {
	if node, ok := h.NodeMap[val]; ok {
		node.Count++
	} else {
		newNode := Node{val, 1}
		h.Data = append(h.Data, &newNode)
		h.SiftUp(len(h.Data) - 1)
		h.NodeMap[val] = &newNode
	}
	h.TotalCount++
}

// Pop the top value of the heap
// put the last node to the first
// and sift down to rebalance
func (h *Heap) Pop() *Node {
	if len(h.Data) <= 1 {
		return nil
	}
	ret := h.Data[1]
	h.Data[1] = h.Data[len(h.Data)-1]
	h.Data = h.Data[:(len(h.Data) - 1)]
	h.SiftDown(1)
	h.TotalCount -= ret.Count
	delete(h.NodeMap, ret.Val)
	return ret
}

// Comparer build compare function depends on
// the heap type
func (h *Heap) Comparer() func(lh, rh *Node) bool {
	if h.Type == MaxHeap {
		return func(lh, rh *Node) bool {
			return lh.Val > rh.Val
		}
	}

	return func(lh, rh *Node) bool {
		return lh.Val < rh.Val
	}
}

// Triple pick is to pick up the min/max value between
// parent node, left child and right child
func (h *Heap) TriplePick(idx int) int {
	lIdx := idx * 2
	rIdx := idx*2 + 1

	if lIdx >= len(h.Data) {
		return idx
	}

	compFunc := h.Comparer()

	if rIdx == len(h.Data) {
		if compFunc(h.Data[lIdx], h.Data[idx]) {
			return lIdx
		}
		return idx
	}

	if compFunc(h.Data[rIdx], h.Data[lIdx]) {
		if compFunc(h.Data[rIdx], h.Data[idx]) {
			return rIdx
		}
		return idx
	}
	if compFunc(h.Data[lIdx], h.Data[idx]) {
		return lIdx
	}
	return idx
}

// SiftDown is to sift the target node down to make
// heap property to be reestablished
func (h *Heap) SiftDown(idx int) {
	maxIdx := h.TriplePick(idx)
	if maxIdx != idx {
		h.Data[maxIdx], h.Data[idx] = h.Data[idx], h.Data[maxIdx]
		h.SiftDown(maxIdx)
	}
}

// SiftUp is to sift the target node up to make
// heap property to be reestablished
func (h *Heap) SiftUp(idx int) {
	if idx <= 1 {
		return
	}
	maxIdx := h.TriplePick(idx / 2)
	if maxIdx != idx/2 {
		h.Data[maxIdx], h.Data[idx/2] = h.Data[idx/2], h.Data[maxIdx]
		h.SiftUp(idx / 2)
	}
}

// TopVal peek the root value for the heap
func (h *Heap) TopVal() uint64 {
	if len(h.Data) <= 1 {
		return 0
	}
	return h.Data[1].Val
}

// TopCount get the count for the heap's top node
func (h *Heap) TopCount() uint64 {
	if len(h.Data) <= 1 {
		return 0
	}
	return h.Data[1].Count
}

// TotalCount get the total count for heap
func (h *Heap) GetTotalCount() uint64 {
	return h.TotalCount
}

// NodeSize get the node count for the heap
func (h *Heap) NodeSize() int {
	return len(h.Data) - 1
}
