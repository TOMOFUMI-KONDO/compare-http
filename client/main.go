package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
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

	// create client with certificate
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := client.Get("https://localhost:44300")
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
