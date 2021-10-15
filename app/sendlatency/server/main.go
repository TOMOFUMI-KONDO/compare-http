package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/TOMOFUMI-KONDO/compare-http/sendlatency/httpversion"

	"github.com/lucas-clemente/quic-go/http3"
)

var (
	crt = "../tls/server.crt"
	key = "../tls/server.key"

	port    string
	version httpversion.Version

	debug bool
)

func init() {
	flag.StringVar(&port, "p", ":443", "server port")
	v := *flag.Int("v", 1, "http version")
	switch v {
	case 1:
		version = httpversion.Ver1
	case 2:
		version = httpversion.Ver2
	case 3:
		version = httpversion.Ver3
	default:
		version = httpversion.Ver1
	}
	flag.BoolVar(&debug, "d", false, "enable debug")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", handler)

	var err error
	switch version {
	case httpversion.Ver1, httpversion.Ver2:
		log.Println("listening http on https://localhost" + port)
		err = http.ListenAndServeTLS(port, crt, key, nil)
	case httpversion.Ver3:
		log.Println("listening http3 on https://localhost" + port)
		err = http3.ListenAndServeQUIC(port, crt, key, nil)
	default:
		log.Fatalf("invalid http version: %d; choose 1 to 3\n", version)
	}

	if err != nil {
		log.Fatalln(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if debug {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(string(dump))
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
	elapsed := time.Now().UnixNano() - sendUnixTime
	fmt.Printf("send latency is %d[ms]\n", elapsed/1000)
}
