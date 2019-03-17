package heap

import (
	"testing"
)

func CreateHeap() *Heap {
	var h = &Heap{
		make([]*Node, 1, 10),
		0,
		MaxHeap,
		make(map[uint64]*Node),
	}

	h.Insert(10)
	h.Insert(9)
	h.Insert(11)
	h.Insert(9)
	return h
}

func TestInsert(t *testing.T) {
	h := CreateHeap()
	if h.TopVal() != 11 {
		t.Fatalf("Top value should be 11 get %d", h.TopVal())
	}
	if h.GetTotalCount() != 4 {
		t.Fatalf("TotalCount should be 4 get %d", h.GetTotalCount())
	}
	if h.NodeSize() != 3 {
		t.Fatalf("NodeSize should be 3 get %d", h.NodeSize())
	}
}

func TestPop(t *testing.T) {
	h := CreateHeap()
	val := h.Pop()
	if val.Val != 11 {
		t.Fatalf("Pop value expect 11 get %d", val.Val)
	}
	if h.GetTotalCount() != 3 {
		t.Fatalf("TotalCount after pop expect 3 get %d", h.GetTotalCount())
	}

	if h.NodeSize() != 2 {
		t.Fatalf("NodeSize after pop expect 2 get %d", h.NodeSize())
	}
	val = h.Pop()
	if val.Val != 10 {
		t.Fatalf("Second Pop value expect 10 get %d", val.Val)
	}
	val = h.Pop()
	if val.Val != 9 || val.Count != 2 {
		t.Fatalf("Third Pop value expect 9 get %d Count should be 2 get %d", val.Val, val.Count)
	}
	val = h.Pop()
	if val != nil {
		t.Fatalf("Pop for empty heap expect nil get %v", val)
	}
}
