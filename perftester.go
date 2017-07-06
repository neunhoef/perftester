package main

import (
  "flag"
  "fmt"
  "io"
  "net/http"
	"os"
	"sort"
  "strconv"
  "time"
)

var blockSize = 1000
var port = 8000
var portStr = "8000"

func hello(w http.ResponseWriter, r *http.Request) {
  io.WriteString(w,"Hello world! Your address is: " + r.RemoteAddr)
}

type Int64Slice []int64

func (s Int64Slice) Len() int {
	return len(s)
}

func (s Int64Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Int64Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func main() {
  flag.IntVar(&blockSize, "blocksize", blockSize, "number of tries in a block")
  flag.IntVar(&port, "port", port, "port to use")
  flag.Parse()
	portStr = strconv.FormatInt(int64(port), 10)

  http.HandleFunc("/", hello);
  go http.ListenAndServe(":" + portStr, nil)

  httpTimes := make([]int64, 0, blockSize)
	fileTimes := make([]int64, 0, blockSize)
	errorCount := 0
	count := 0

  transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		DisableKeepAlives:  false,
	}
  client := &http.Client{Transport: transport};

  for {
    time.Sleep(100000000)

		// Measure a HTTP roundtrip:
		start := time.Now()
		req, _ := http.NewRequest("GET", "http://localhost:" + portStr, nil)
		req.Header = map[string][]string{"Keep-Alive": {"300"}}
		r, e := client.Do(req)
		end := time.Now()
		if e != nil {
		  errorCount += 1
		} else {
		  r.Body.Close()
		}
		httpTimes = append(httpTimes, int64(end.Sub(start)) / 1000)

    // Measure a file write:
		start = time.Now()
		name := "/tmp/testfile" + strconv.FormatInt(int64(count), 10)
		f, e := os.Create(name)
		if e != nil {
		  errorCount += 1
		} else {
		  msg := []byte("Hallo")
		  _, e := f.Write(msg)
			if e != nil {
			  errorCount += 1
			}
			if f.Close() != nil {
			  errorCount += 1
			}
		}
		_ = os.Remove(name)
		end = time.Now()
		fileTimes = append(fileTimes, int64(end.Sub(start)) / 1000)

		if len(httpTimes) >= blockSize {
		  fmt.Println(end.Format(time.RFC3339))
      sort.Sort(Int64Slice(httpTimes))
			var sum int64 = 0
			for _, t := range(httpTimes) {
			  sum += t
			}
			fmt.Println("Http times: avg=", int(float64(sum) / float64(blockSize)),
			            " 50%=", httpTimes[blockSize/2],
									" 90%=", httpTimes[int(float64(blockSize) * 0.9)],
									" 99%=", httpTimes[int(float64(blockSize) * 0.99)])
			fmt.Printf("  10 longest times:")
			first := blockSize - 10
			if first < 0 {
			  first = 0
			}
			for i := first; i < blockSize; i++ {
			  fmt.Printf(" %d", httpTimes[i])
			}
			fmt.Println()

      sort.Sort(Int64Slice(fileTimes))
			sum = 0
			for _, t := range(fileTimes) {
			  sum += t
			}
			fmt.Println("File times: avg=", int(float64(sum) / float64(blockSize)),
			            " 50%=", fileTimes[blockSize/2],
									" 90%=", fileTimes[int(float64(blockSize) * 0.9)],
									" 99%=", fileTimes[int(float64(blockSize) * 0.99)])
			fmt.Printf("  10 longest times:")
			for i := first; i < blockSize; i++ {
			  fmt.Printf(" %d", fileTimes[i])
			}
			fmt.Println()

			fmt.Println("Cumulated error count:", errorCount, "\n")
			httpTimes = httpTimes[:0]
			fileTimes = fileTimes[:0]
		}
  }
}
