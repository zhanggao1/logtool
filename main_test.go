package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"math/rand"
	"testing"
)

func Benchmark_Analyse(b *testing.B) {
	file, err := ioutil.TempFile("", "benchmark")
	if err != nil {
		b.Fatal(err)
	}
	writeBuf := bufio.NewWriter(file)
	var buf = make([]byte, 8)
	for i := 0; i < 10*1000*1000; i++ {
		binary.BigEndian.PutUint64(buf, uint64(rand.Int31()))
		writeBuf.Write(buf)
	}
	writeBuf.Flush()
	file.Close()
	Analyse(10*1000*1000, file.Name(), []float32{0.1, 0.05, 0.01})
}

func Test_BuildHeap(t *testing.T) {
	var data = make([]byte, 2048, 2048)
	for j := 0; j < 2; j++ {
		for i := 0; i < 1024/8; i++ {
			binary.BigEndian.PutUint64(data[j*1024+8*i:], uint64(i))
		}
	}
	reader := bytes.NewReader([]byte(data))
	BuildHeap(16, reader)
	top := topNHeap.Pop()
	if top.Val != 1024/8-1-16/2+1 {
		t.Fatalf("Top value expect 120 after build heap get %d", top.Val)
	}
	if topNHeap.GetTotalCount() != 16-2 {

		t.Fatalf("Heap total count expect 30 after build heap get %d", topNHeap.GetTotalCount())
	}
}
