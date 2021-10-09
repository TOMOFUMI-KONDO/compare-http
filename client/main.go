package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

	"golang.org/x/net/http2"
)

func main() {
	// load certificate
	cert, err := ioutil.ReadFile("tls/ca.crt")
	if err != nil {
		log.Fatalln(err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)
	//tlsConfig := &tls.Config{RootCAs: certPool}
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	tlsConfig.BuildNameToCertificate()

	version := flag.Int("version", 3, "http version")
	flag.Parse()

	// create client with certificate
	var client *http.Client
	switch *version {
	case 3:
		log.Fatalln("not implemented")
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
		log.Fatalf("Inavlid version: %d\n", *version)
	}

	resp, err := client.Get("https://localhost")
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
