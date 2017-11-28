package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"bufio"
 	"os"
	"strconv"
	"net/http/httputil"
	"time"
)

type urlCounter struct{
	url	string
	counter int
}

type result struct {
	Foo string
}

type Test struct {
	url    string
	result interface{}
}

var badResponseCode []string
var badStringSize []string
var counterSlice []urlCounter
var incompleteRequest []string

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter filePath of textfile: ")
	//C:/Users/byang8/go/src/test/test.txt

	filePath, _ := reader.ReadString('\n')

	filePath =  filePath[:len(filePath)-len("11")]

	tempTests := readLines(filePath)
	for x := range tempTests {
		counterSlice = append(counterSlice, urlCounter{url: tempTests[x].url, counter: 0})
	}

	fmt.Print("\nEnter amount of iterations to pass through all services in textfile: ")

  IterationsTemp, err1 := reader.ReadString('\r')
	if err1 != nil {
		fmt.Println(err1)
	}

	Iterations :=  IterationsTemp[:len(IterationsTemp)-len("1")]

	intIterationsReg, err3 := strconv.ParseInt(Iterations, 10, 0)
	if err3 != nil {
		fmt.Println(err3)
	}
	intIterations := int64(intIterationsReg)

	reader.ReadString('\n')

	fmt.Print("\nEnter amount of seconds before terminating all threads\n(recommended no less than 60 for one pass of textfile): ")
  SecondsTemp, err2 := reader.ReadString('\r')
	if err2 != nil {
		fmt.Println(err2)
	}

	reader.ReadString('\n')

	Seconds :=  SecondsTemp[:len(SecondsTemp)-len("1")]

	intSeconds, err4 := strconv.ParseInt(Seconds, 10, 0)
	if err4 != nil {
		fmt.Println(err4)
	}


	var wg sync.WaitGroup

	var j int64
	for j = 0; j < intIterations; j++ {
		for i, test := range tempTests {
			fmt.Println(test.url)
			wg.Add(1)
			//err := getJSON(test.url, &test.result)
			go getJSON(wg, j, i, test.url, &test.result, intIterations)
		}
	}

	time.Sleep(time.Duration(intSeconds * 1000) * time.Millisecond)

	fmt.Printf("\n\n\nbad response code urls:\n")
	for i := range badResponseCode {
		fmt.Printf(badResponseCode[i] + "\n")
	}
	fmt.Printf("\nbad string size urls (less than 1MB):\n")
	for i := range badStringSize {
		fmt.Printf(badStringSize[i] + "\n")
	}
	fmt.Printf("\nincomplete requests by the time limit:\n")
	for i := range incompleteRequest {
		fmt.Printf(incompleteRequest[i] + "\n")
	}
	fmt.Printf("\ntotal number of iterations of an url completed\n")
	for i := range counterSlice {
		fmt.Printf("\nfor " + counterSlice[i].url + " :: " + strconv.Itoa(counterSlice[i].counter) + " out of " + strconv.FormatInt(intIterations, 10) + " completed\n")
	}
}

func readLines(path string) ([]Test) {
	var tempTempTests []Test
  file, err := os.Open(path)
  if err != nil {
    return nil
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
	for i := range lines {
		tempTempTests = append(tempTempTests, Test{url: lines[i], result: new(result)})
	}
  return tempTempTests
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func getJSON(wg sync.WaitGroup, j int64, i int, url string, result interface{}, intIterations int64) error {

	for q := range counterSlice {
		if counterSlice[q].url == url {
			counterSlice[q].counter++
			break
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("cannot fetch URL %q: %v", url, err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 503 {
			if (!contains(incompleteRequest, url)) {
				incompleteRequest = append(incompleteRequest, url)
			}
		} else {
			if (!contains(badResponseCode, url)) {
				badResponseCode = append(badResponseCode, url + " :: " + strconv.FormatInt(int64(resp.StatusCode), 10))
			}
		}
		fmt.Errorf("unexpected http GET status: %s", resp.Status)
		return nil
	}
	fmt.Printf("\nresponse code %d out of " + strconv.FormatInt(intIterations, 10) + " of %s : %d\n", (j + 1), url, resp.StatusCode)
	dump, err := httputil.DumpResponse(resp, true)
	fmt.Printf("\ncontent length %d out of " + strconv.FormatInt(intIterations, 10) + " of %s : %d\n", (j + 1), url, len(dump))
	if (len(dump) < 100000) {
		if (!contains(badStringSize, url)) {
			badStringSize = append(badStringSize, url + " :: " + strconv.FormatInt(int64(len(dump)), 10))
		}
	}
	//fmt.Printf("\ncontent length %d : %d\n", i, resp.ContentLength)
	//err = json.Unmarshal(resp.Body, result)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		fmt.Errorf("cannot decode JSON: %v", err)
		return nil
	}
	if err != nil {
		fmt.Printf("\ntest %d: error %v\n", i, err)
	} else {
		fmt.Printf("\ntest %d out of " + strconv.FormatInt(intIterations, 10) + " of %s: ok with result: \n%#v\n", (j + 1), url, result)
	}

	defer wg.Done()

	return nil
}
