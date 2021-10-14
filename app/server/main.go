package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/lucas-clemente/quic-go/http3"
)

var (
	crt     = "tls/server.crt"
	key     = "tls/server.key"
	port    string
	version int
)

type Body struct {
	File string `json:"file"`
}

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
	// dump request
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Printf("[ERROR] Failed to DumpRequest.\n%v\n", err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))

	// read body
	bodyByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body.\n%v\n", err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	var body Body
	if err := json.Unmarshal(bodyByte, &body); err != nil {
		log.Printf("[ERROR] Failed to unmarshal json.\n%v\n", err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	// load file
	file, err := os.Open("server/assets/" + body.File)
	if err != nil {
		log.Printf("[ERROR] Failed to open file.\n%v\n", err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// return response
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("[ERROR] Failed to copy response.\n%v\n", err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
}
