package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/lucas-clemente/quic-go/http3"
)

var (
	crt = "tls/server.crt"
	key = "tls/server.key"
)

func main() {
	port := ":44300"

	http.HandleFunc("/", handler)

	version := flag.Int("version", 3, "http version")
	flag.Parse()

	var err error
	switch *version {
	case 3:
		log.Println("listening http3 on https://localhost" + port)
		err = http3.ListenAndServeQUIC(port, crt, key, nil)
	case 2:
		log.Println("listening http2 on https://localhost" + port)
		err = http.ListenAndServeTLS(port, crt, key, nil)
	}

	if err != nil {
		log.Fatalln(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w, "ok")
}
