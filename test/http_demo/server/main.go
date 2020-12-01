package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	ADDR = ":8080"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/abc", handlerFunc)
	// 创建服务器
	server := &http.Server{
		Addr:         ADDR,
		Handler:      mux,
		WriteTimeout: time.Second * 5,
	}
	fmt.Printf("Listen and server on %v\n", ADDR)
	log.Fatal(server.ListenAndServe())
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(r.Header)
	io.WriteString(w, fmt.Sprintf("header: %s\n", b))
	io.WriteString(w, fmt.Sprintf("time: %s\n", time.Now().Format(time.RFC3339)))
}
