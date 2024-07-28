package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handleRequest)
	server := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		log.Printf("Request took %v", time.Since(startTime))
	}()

	// 记录请求日志
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	// 转发请求到服务端
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("http://your-service-endpoint")
	if err != nil {
		log.Printf("Error forwarding request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 将响应转发回客户端
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(body)
}
