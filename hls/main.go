package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port          = "8000"
	mainPath      = "/"
	storagePath   = "/storage/"
	playerServer  = http.FileServer(http.Dir("hls/player"))
	storageServer = http.FileServer(http.Dir("hls/storage"))
)

func main() {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:           "localhost:" + port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	mux.Handle(mainPath, addHeaders(http.StripPrefix(mainPath, playerServer)))
	mux.Handle(storagePath, addHeaders(http.StripPrefix(storagePath, storageServer)))

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	fmt.Println("HLS Player App started!")
	log.Printf("Starting HTTP server on %s", port)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped
	fmt.Println("HLS Player App stopped")

}

// addHeaders Функция добавляет поддержку CORS
func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}
