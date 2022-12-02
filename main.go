package main

import (
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	protocol   = "unix"
	socketPath = "/run/earth/web.sock"
)

func helloEarth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Earth❤"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloEarth)

	// listener作成
	listener, err := net.Listen(protocol, socketPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("error: %v", err)
		}
	}()

	// ソケットはデフォルトだと同じユーザーしか使えないのでNginxが使えるようにしておく
	err = os.Chmod(socketPath, fs.ModePerm)
	if err != nil {
		log.Printf("error: %v", err)
	}

	// graceful shutdown
	go func() {
		sig := make(chan os.Signal, 2)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		if err := listener.Close(); err != nil {
			log.Printf("error: %v", err)
		}
		os.Exit(1)
	}()

	log.Println("Mother Earth is carefully Listening...")
	if err := http.Serve(listener, mux); err != nil {
		log.Fatalf("error: %v", err)
	}
}
