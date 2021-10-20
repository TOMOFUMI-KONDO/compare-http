package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"

	"github.com/TOMOFUMI-KONDO/compare-http/continuoussend/chart"

	"github.com/lucas-clemente/quic-go/http3"
)

var (
	crt = "../tls/server.crt"
	key = "../tls/server.key"

	port    string
	version int
	debug   bool
)

func init() {
	flag.StringVar(&port, "port", ":443", "server port")
	flag.IntVar(&version, "version", 1, "http version")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", handle)
	http.HandleFunc("/fin", handleFin)
	fs := http.FileServer(http.Dir("out"))
	http.Handle("/charts/", http.StripPrefix("/charts/", fs))

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

func handle(w http.ResponseWriter, r *http.Request) {
	if debug {
		if err := dumpRequest(r); err != nil {
			fmt.Fprintf(os.Stderr, "failed to dump request; %s\n", err)
		}
	}

	sendAt, err := strconv.ParseInt(r.FormValue("send_at"), 10, 64)
	if err != nil {
		fmt.Printf("failed to parse send_at; %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	latencyNanoSec := time.Now().UnixNano() - sendAt
	latencyMilSec := int(latencyNanoSec / int64(math.Pow(10, 6)))

	chart.Record(latencyMilSec)
}

func handleFin(w http.ResponseWriter, r *http.Request) {
	if debug {
		if err := dumpRequest(r); err != nil {
			fmt.Fprintf(os.Stderr, "failed to dump request; %s\n", err)
		}
	}

	if err := chart.Render(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to render chart; %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func dumpRequest(r *http.Request) error {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}

	fmt.Println("\n===============request===============")
	fmt.Println(string(dump))
	fmt.Println("=====================================")

	fmt.Printf("size: %d[byte]\n", len(dump))

	return nil
}
