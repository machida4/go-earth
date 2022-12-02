package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	protocol = "unix"
	address  = "/run/earth/web.sock"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Earth❤"))
	})

	listener, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("error: %v", err)
		}
	}()

	shutdown(listener)
	if err := http.Serve(listener, mux); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func shutdown(listener net.Listener) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		//s := <-c
		if err := listener.Close(); err != nil {
			log.Printf("error: %v", err)
		}
		os.Exit(1)
	}()
}
