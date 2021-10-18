package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/lucas-clemente/quic-go/http3"

	"golang.org/x/net/http2"
)

var (
	host, port string
	version    int
	files      = []string{"1M.txt", "10M.txt", "100M.txt", "1000M.txt"}
	times      = 3 // try this times, and take average.
)

func init() {
	flag.StringVar(&host, "h", "localhost", "server host")
	flag.StringVar(&port, "p", ":443", "server port")
	flag.IntVar(&version, "v", 3, "http version")
	flag.Parse()
}

func main() {
	// FIXME: This is ignoring server certificate. Don't use in production.
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	var client *http.Client
	switch version {
	case 3:
		client = &http.Client{
			Transport: &http3.RoundTripper{
				TLSClientConfig: tlsConfig,
			},
		}
	case 2:
		client = &http.Client{
			Transport: &http2.Transport{
				TLSClientConfig: tlsConfig,
			},
		}
	case 1:
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}
	default:
		log.Fatalf("Inavlid version: %d; choole 1 to 3\n", version)
	}

	for _, f := range files {
		fmt.Printf("file: %s\n", f)

		var sum int64 = 0

		for i := 0; i < times+1; i++ {
			start := time.Now()

			request(client, f)
			// to ignore overhead of initializing connection
			if i == 0 {
				continue
			}

			elapsed := time.Since(start).Milliseconds()
			fmt.Printf("response time: %d[ms]\n", elapsed)
			sum = sum + elapsed
		}

		average := sum / int64(times)
		fmt.Printf("average: %d[ms]\n\n", average)
	}
}

func request(client *http.Client, file string) {
	query := url.Values{
		"file": {file},
	}

	resp, err := client.Get("https://" + host + port + "?" + query.Encode())
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		log.Fatalln(err)
	}
}
