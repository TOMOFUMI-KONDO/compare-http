package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/lucas-clemente/quic-go/http3"
)

var (
	crt = "../tls/server.crt"
	key = "../tls/server.key"

	port    string
	version int
	debug   bool

	latencies []int64
)

func init() {
	flag.StringVar(&port, "p", ":443", "server port")
	flag.IntVar(&version, "v", 1, "http version")
	flag.BoolVar(&debug, "d", false, "enable debug")
	flag.Parse()
}

func main() {
	http.HandleFunc("/record", handleRecord)
	http.HandleFunc("/stat", handleStat)

	var err error
	switch version {
	case 1, 2:
		log.Println("listening http on https://localhost" + port)
		err = http.ListenAndServeTLS(port, crt, key, nil)
	case 3:
		log.Println("listening http3 on https://localhost" + port)
		err = http3.ListenAndServeQUIC(port, crt, key, nil)
	default:
		log.Fatalf("invalid http version: %d; choose 1 to 3\n", version)
	}

	if err != nil {
		log.Fatalln(err)
	}
}

func handleRecord(w http.ResponseWriter, r *http.Request) {
	if debug {
		dumpRequest(r)
	}

	//read and discard load file
	file, _, err := r.FormFile("load")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	if _, err := io.ReadAll(file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// measure send latency
	sendUnixTime, err := strconv.ParseInt(r.FormValue("send_time_unix_nano"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	latency := time.Now().UnixNano() - sendUnixTime
	latencyMs := latency / int64(math.Pow(10, 6))

	latencies = append(latencies, latencyMs)

	fmt.Printf("send latency is %d[ms]\n", latencyMs)
}

func handleStat(_ http.ResponseWriter, r *http.Request) {
	if debug {
		dumpRequest(r)
	}

	var sum int64 = 0
	for _, v := range latencies {
		sum = sum + v
	}
	average := sum / int64(len(latencies))

	fmt.Printf("send latency average: %d[ms]\n\n", average)

	// reset stat
	latencies = nil
}

func dumpRequest(r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(dump))
}
