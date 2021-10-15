package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"

	"github.com/lucas-clemente/quic-go/http3"
	"golang.org/x/net/http2"
)

var (
	host, port string
	version    int
	debug      bool
	times      int // try this times, and take
	file       string
)

func init() {
	flag.StringVar(&host, "h", "localhost", "server hostname")
	flag.StringVar(&port, "p", ":443", "server port")
	flag.IntVar(&version, "v", 1, "http version")
	flag.BoolVar(&debug, "d", false, "enable debug")
	flag.IntVar(&times, "t", 10, "try times. take average of these")
	flag.StringVar(&file, "f", "1K.txt", "sent file")
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

	for i := 0; i < times; i++ {
		send(client, "client/assets/"+file)
	}

	resp, err := client.Get("https://" + host + port + "/stat")
	if err != nil {
		log.Fatalln(err)
	}

	if debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(dump))
	} else {
		// read and discard response body
		if _, err := io.ReadAll(resp.Body); err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
	}
}

func send(client *http.Client, filepath string) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// set load file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	fileWriter, err := writer.CreateFormFile("load", file.Name())
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := io.Copy(fileWriter, file); err != nil {
		log.Fatalln(err)
	}

	// set send time
	now := time.Now().UnixNano()
	if err := writer.WriteField("send_time_unix_nano", strconv.FormatInt(now, 10)); err != nil {
		log.Fatalln(err)
	}

	// write all
	if err := writer.Close(); err != nil {
		log.Fatalln(err)
	}

	// send request
	resp, err := client.Post("https://"+host+port+"/record", writer.FormDataContentType(), &buffer)
	if err != nil {
		log.Fatalln(err)
	}

	if debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(dump))
	} else {
		// read and discard response body
		if _, err := io.ReadAll(resp.Body); err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
	}
}
