package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/lucas-clemente/quic-go/http3"

	"golang.org/x/net/http2"
)

var (
	host, port string
	version    int
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

	resp, err := client.Get("https://" + host + port)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	fmt.Printf("Protocol Version:%s\n\n", resp.Proto)
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(dump))
}
