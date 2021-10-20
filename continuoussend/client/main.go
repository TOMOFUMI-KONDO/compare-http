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
	"strconv"
	"time"

	"github.com/lucas-clemente/quic-go/http3"
	"golang.org/x/net/http2"
)

var (
	host, port            string
	version, tick, period int
	debug                 bool
)

func init() {
	flag.StringVar(&host, "host", "localhost", "server hostname")
	flag.StringVar(&port, "port", ":443", "server port")
	flag.IntVar(&version, "version", 1, "http version")
	flag.IntVar(&tick, "tick", 10, "tick interval [ms]")
	flag.IntVar(&period, "period", 60, "continue to send units for this period [s]")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.Parse()
}

func main() {
	client, err := makeClient(version)
	if err != nil {
		log.Fatalf("failed to make client; %s\n", err.Error())
	}

	ticker := time.NewTicker(time.Duration(tick) * time.Millisecond)
	done := make(chan bool)

	fmt.Printf("\n=====send start=====\n")

	go func() {
		for {
			select {
			case <-ticker.C:
				if err = send(client); err != nil {
					log.Printf("failed to send request; %s\n", err.Error())
				}
			case <-done:
				return
			}
		}
	}()

	time.Sleep(time.Duration(period) * time.Second)
	ticker.Stop()
	done <- true

	if err = fin(client); err != nil {
		log.Fatalf("failed to send fin; %s\n", err)
	}

	fmt.Printf("\n=====send finish=====\n")
}

func makeClient(version int) (*http.Client, error) {
	// FIXME: This is ignoring server certificate. Don't use in production.
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	switch version {
	case 3:
		return &http.Client{
			Transport: &http3.RoundTripper{
				TLSClientConfig: tlsConfig,
			},
		}, nil
	case 2:
		return &http.Client{
			Transport: &http2.Transport{
				TLSClientConfig: tlsConfig,
			},
		}, nil
	case 1:
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}, nil
	default:
		return nil, fmt.Errorf("Inavlid version: %d; choole 1 to 3\n", version)
	}
}

func send(client *http.Client) error {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	now := time.Now().UnixNano()
	if err := writer.WriteField("send_at", strconv.FormatInt(now, 10)); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	resp, err := client.Post("https://"+host+port, writer.FormDataContentType(), &buffer)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if debug {
		dump, _ := httputil.DumpResponse(resp, false)
		fmt.Printf("%s\n%s\n", dump, body)
	}

	return nil
}

func fin(client *http.Client) error {
	resp, err := client.Post("https://"+host+port+"/fin", "application/x-www-url-encoded", nil)
	if err != nil {
		return err
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if debug {
		dump, _ := httputil.DumpResponse(resp, false)
		fmt.Println(string(dump))
	}

	return nil
}
