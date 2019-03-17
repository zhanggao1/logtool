package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestBuildHeap(t *testing.T) {
	var data = make([]byte, 2048, 2048)
	for j := 0; j < 2; j++ {
		for i := 0; i < 1024/8; i++ {
			binary.BigEndian.PutUint64(data[j*1024+8*i:], uint64(i))
		}
	}
	reader := bytes.NewReader([]byte(data))
	BuildHeap(16, reader)
	fmt.Println(topNHeap.GetTotalCount())
	top := topNHeap.Pop()
	if top.Val != 1024/8-1-16/2+1 {
		t.Fatalf("Top value expect 120 after build heap get %d", top.Val)
	}
	if topNHeap.GetTotalCount() != 16-2 {

		t.Fatalf("Heap total count expect 30 after build heap get %d", topNHeap.GetTotalCount())
	}
}
