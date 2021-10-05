package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

var (
	port = ":8000"
)

func main(){
	http.HandleFunc("/", handler)
	log.Println("listening localhost" + port)
	if err:= http.ListenAndServe(port,nil);err!=nil{
		log.Fatalln(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	dump,err:=httputil.DumpRequest(r,true)
	if err!=nil{
		http.Error(w,fmt.Sprint(err),http.StatusInternalServerError)
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w,"<html><body>hello</body></html>")
}
