package main

import (
	"fmt"
	"sort"
	"strconv"
)

var referenceResult = `29568666068035183841425683795340791879727309630931025356555_4958044192186797981418233587017209679042592862002427381542`

func main() {
	TestExecutePipeline()
}

func TestExecutePipeline() {
	var data []int = []int{0, 1}
	hashes := make([]string, len(data))

	for i := range data {
		singleHashResult := TestSingleHash(i, data[i])
		multiHashResult := TestMultiHash(singleHashResult)
		hashes[i] = multiHashResult
		fmt.Println()
	}

	finalResult := TestCombineResults(hashes...)

	// CombineResults 29568666068035183841425683795340791879727309630931025356555_4958044192186797981418233587017209679042592862002427381542
	fmt.Printf("CombineResults %s\n", finalResult)

	if finalResult == referenceResult {
		fmt.Println("Good job!")
	} else {
		fmt.Println("Something wrong")
	}
}

func TestSingleHash(iter, data int) string {
	// 0 SingleHash data 0
	// 0 SingleHash md5(data) cfcd208495d565ef66e7dff9f98764da
	// 0 SingleHash crc32(md5(data)) 502633748
	// 0 SingleHash crc32(data) 4108050209
	// 0 SingleHash result 4108050209~502633748

	dataStr := strconv.Itoa(data)
	fmt.Printf("%d SingleHash data %s\n", iter, dataStr)

	hashMd5 := DataSignerMd5(dataStr)
	fmt.Printf("%d SingleHash md5(data) %s\n", iter, hashMd5)

	secondCrc32 := DataSignerCrc32(hashMd5)
	fmt.Printf("%d SingleHash crc32(md5(data)) %s\n", iter, secondCrc32)

	firstCrc32 := DataSignerCrc32(dataStr)
	fmt.Printf("%d SingleHash crc32(data) %s\n", iter, firstCrc32)

	result := fmt.Sprintf("%s~%s", firstCrc32, secondCrc32)
	fmt.Printf("%d SingleHash result %s\n", iter, result)

	return result
}

func TestMultiHash(data string) string {
	// 4108050209~502633748 MultiHash: crc32(th+step1)) 0 2956866606
	// 4108050209~502633748 MultiHash: crc32(th+step1)) 1 803518384
	// 4108050209~502633748 MultiHash: crc32(th+step1)) 2 1425683795
	// 4108050209~502633748 MultiHash: crc32(th+step1)) 3 3407918797
	// 4108050209~502633748 MultiHash: crc32(th+step1)) 4 2730963093
	// 4108050209~502633748 MultiHash: crc32(th+step1)) 5 1025356555
	// 4108050209~502633748 MultiHash result: 29568666068035183841425683795340791879727309630931025356555

	var result string

	for i := 0; i < 6; i++ {
		newData := strconv.Itoa(i) + data
		hash := DataSignerCrc32(newData)
		fmt.Printf("%s MultiHash: crc32(th+step1)) %d %s\n", data, i, hash)

		result += hash
	}

	fmt.Printf("%s MultiHash result: %s\n", data, result)
	return result
}

func TestCombineResults(hashes ...string) string {
	var result string

	sort.Strings(hashes)

	for _, hash := range hashes {
		result += hash + "_"
	}

	result = result[:len(result)-1]
	return result
}

func ExecutePipeline(jobs ...job) {
	// in := make(chan interface{})
	// out := make(chan interface{})
}

var SingleHash = func(in, out chan interface{}) {

}

var MultiHash = func(in, out chan interface{}) {

}

var CombineResults = func(in, out chan interface{}) {

}
