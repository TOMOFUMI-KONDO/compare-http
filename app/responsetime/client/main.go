package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/lucas-clemente/quic-go/http3"

	"golang.org/x/net/http2"
)

var (
	host, port, file string
	version          int
)

func init() {
	flag.StringVar(&host, "h", "localhost", "server host")
	flag.StringVar(&port, "p", ":443", "server port")
	flag.IntVar(&version, "v", 3, "http version")
	flag.StringVar(&file, "f", "1M.txt", "request file name")
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

	query := url.Values{
		"file": {file},
	}

	start := time.Now()

	resp, err := client.Get("https://" + host + port + "?" + query.Encode())
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	responseTime := time.Now().Sub(start)
	fmt.Printf("response_time: %s\n\n", responseTime.String())

	dump, err := httputil.DumpResponse(resp, false)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(dump))
}
