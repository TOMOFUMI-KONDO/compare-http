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

	"github.com/TOMOFUMI-KONDO/compare-http/sendlatency/httpversion"

	"github.com/lucas-clemente/quic-go/http3"
	"golang.org/x/net/http2"
)

var (
	host, port string
	version    httpversion.Version
	debug      bool

	assetsBasePath = "client/assets/"
)

func init() {
	flag.StringVar(&host, "h", "localhost", "server hostname")
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
	// FIXME: This is ignoring server certificate. Don't use in production.
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	var client *http.Client
	switch version {
	case httpversion.Ver3:
		client = &http.Client{
			Transport: &http3.RoundTripper{
				TLSClientConfig: tlsConfig,
			},
		}
	case httpversion.Ver2:
		client = &http.Client{
			Transport: &http2.Transport{
				TLSClientConfig: tlsConfig,
			},
		}
	case httpversion.Ver1:
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}
	default:
		log.Fatalf("Inavlid version: %d; choole 1 to 3\n", version)
	}

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// set load file
	file, err := os.Open(assetsBasePath + "1M.txt")
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
	resp, err := client.Post("https://"+host+port, writer.FormDataContentType(), &buffer)
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
