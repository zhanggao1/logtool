// Package main is the main entry
// for the whole project
// it contains function for read all log files, build up heap .etc
package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"heap"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Temporary file name to maintain all valid response time
// read from the log file
const TmpFilePrefix string = "tmp"

// String used to check if the log line is valid
const HttpVerbGet string = "GET"

// Http success code, used to check if the log line is valid
const HttpOk string = "200"

const (
	// The word index for the http verb in each line of the log file
	HttpVerbFieldIdx int = 2
	// The http state index for the http verb in each line of the log file
	StateFieldIdx int = 4
	// The response time index for the http verb in each line of the log file
	TimeFieldIdx int = 5
)

// The initial heap capacity for the top n heap
const HeapInitCapacity int = 1024

// Number of uint64 to be fetch in each file read
const BatchReadSize int = 512

// The top p heap used to host the top n response time
var topNHeap = heap.NewHeap(HeapInitCapacity, heap.MinHeap)

// Read from the temporary file which maintain all the valid response time
// Build up and refresh the top n heap
func BuildHeap(size uint64, buf io.Reader) {
	var bufBytes = make([]byte, 8*BatchReadSize)
	for {
		byteSize, err := buf.Read(bufBytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		for i := 0; i < byteSize/8; i++ {
			cost := uint64(binary.BigEndian.Uint64(bufBytes[i*8 : (i+1)*8]))
			if topNHeap.GetTotalCount() < size {
				topNHeap.Insert(cost)
			} else if cost > topNHeap.TopVal() {
				topNHeap.Insert(cost)
				if topNHeap.GetTotalCount()-topNHeap.TopCount() >= size {
					topNHeap.Pop()
				}
			}
		}
	}
}

// Check if the line is valid log of a READ API
func isLineValid(words []string) bool {
	if len(words) < TimeFieldIdx+1 {
		return false
	}
	httpVerb := words[HttpVerbFieldIdx]
	if len(httpVerb) < 4 || httpVerb[1:] != HttpVerbGet {
		return false
	}
	if words[StateFieldIdx] != HttpOk {
		return false
	}
	return true
}

// Open the temporary file to read all the valid response time
// Build up min-heap and feed it with all the valid response time
// Figure out the result
// Resize the heap for every request percentile
func Analyse(totalCount uint64, fileName string, infoRates []float32) {
	timeFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	buf := bufio.NewReader(timeFile)
	topRate := infoRates[0]
	cacheSize := uint64(float32(totalCount) * topRate)

	BuildHeap(cacheSize, buf)

	for _, rate := range infoRates {
		newSize := uint64(rate * float32(totalCount))
		for topNHeap.TotalCount > newSize {
			topNHeap.Pop()
		}
		fmt.Printf("%d%% of requests return a response in %d ms\n", int((1-rate)*100), topNHeap.TopVal())
	}
	timeFile.Close()
	err = os.Remove(fileName)
	if err != nil {
		panic(err)
	}
}

// ScanFile is the main entry for each goroutine to read on log file
// Send the valid response time to the collect goroutine by channel
// Send 0 to indicate the end of file
func ScanFile(name string, analyseChan chan uint64) error {
	logFile, err := os.Open(name)
	if err != nil {
		return err
	}
	defer logFile.Close()
	buf := bufio.NewReader(logFile)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				analyseChan <- 0
				return nil
			}
			return err
		}
		line = strings.TrimSpace(line)
		words := strings.Split(line, " ")
		if isLineValid(words) {
			val, err := strconv.ParseInt(words[TimeFieldIdx], 10, 64)
			if err == nil {
				analyseChan <- uint64(val)
			}
		}
	}
	return nil
}

// parse user input
// parse the percentile list to the rate list
// which will be used to fetch the top biggest
// response time
func parseInput() (path string, rateList []float32) {
	var percentileList string
	flag.StringVar(&path, "path", "/var/log/httpd/", "path of log files")
	flag.StringVar(&percentileList, "percentile-list", "90,95,99", "list for the target percentile response time for the READ API request")
	flag.Parse()

	rateList = make([]float32, 0, 3)
	for _, val := range strings.Split(percentileList, ",") {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			panic("Invalid input for percentile list:" + val)
		}
		if intVal < 0 || intVal > 100 {
			panic("Invalid input for percentile list")
		}
		rateList = append(rateList, float32(100-intVal)/100.0)
	}
	if len(rateList) < 1 {
		panic("percentile list should have at least one value")
	}
	sort.Slice(rateList, func(i, j int) bool { return rateList[i] > rateList[j] })
	return
}

// Start entry of the tool, works as follow:
// 1.read all the valid log files in the specified directory
// 2.write all the valid response time to the temp file
// 3.analyse all the value in the temp file to build up a top x% heap
// 4.shrink the heap to sizes of other specified rate
func main() {
	path, rateList := parseInput()
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	var reflectItems = make([]reflect.SelectCase, 0, 100)
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".log") {
			newFileChan := make(chan uint64, 10)
			go ScanFile(path+"/"+file.Name(), newFileChan)

			reflectItems = append(reflectItems,
				reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(newFileChan),
				})
		}
	}
	fileForCosts, err := ioutil.TempFile(".", TmpFilePrefix)
	if err != nil {
		panic(err)
	}
	writeBuf := bufio.NewWriter(fileForCosts)
	var totalLines uint64 = 0
	var buf = make([]byte, 8)
	for {
		chosen, value, _ := reflect.Select(reflectItems)
		cost := value.Uint()
		if cost == 0 { // read file end
			if chosen == 0 && len(reflectItems) == 1 {
				break
			}
			reflectItems = append(reflectItems[:chosen], reflectItems[chosen+1:]...)
			continue
		}
		totalLines++
		binary.BigEndian.PutUint64(buf, cost)
		writeBuf.Write(buf)
	}
	writeBuf.Flush()
	fileForCosts.Close()
	Analyse(totalLines, fileForCosts.Name(), rateList)
}
