package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func main(){
	port := ":44300"

	http.HandleFunc("/", handler)
	log.Println("listening https://localhost" + port)

	if err:= http.ListenAndServeTLS(port,"tls/server.crt", "tls/server.key",nil);err!=nil{
		log.Fatalln(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	dump,err:=httputil.DumpRequest(r,true)
	if err!=nil{
		http.Error(w,fmt.Sprint(err),http.StatusInternalServerError)
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w,"ok")
}
