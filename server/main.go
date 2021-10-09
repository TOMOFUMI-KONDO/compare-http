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
	crt     = "tls/server.crt"
	key     = "tls/server.key"
	port    string
	version int
)

func init() {
	flag.StringVar(&port, "p", ":443", "server port")
	flag.IntVar(&version, "v", 3, "http version")
	flag.Parse()
}

func main() {

	http.HandleFunc("/", handler)

	var err error
	switch version {
	case 3:
		log.Println("listening http3 on https://localhost" + port)
		err = http3.ListenAndServeQUIC(port, crt, key, nil)
	case 2, 1:
		log.Println("listening http on https://localhost" + port)
		err = http.ListenAndServeTLS(port, crt, key, nil)
	default:
		log.Fatalf("Inavlid version: %d; choole 1 to 3\n", version)
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
