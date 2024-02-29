package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// EchoText Метод получающий и отправляющий обычные текстовые сообщения
func EchoText(ws *websocket.Conn) {
	for {
		var reply string

		if err := websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Printf("Can't receive: %v\n", err)
			break
		}

		if err := websocket.Message.Send(ws, reply); err != nil {
			fmt.Printf("Can't send: %v\n", err)
			break
		}

	}
}

// EchoJSON Метод получающий и отправляющий JSON сообщения
func EchoJSON(ws *websocket.Conn) {
	for {
		var reply Message

		if err := websocket.JSON.Receive(ws, &reply); err != nil {
			fmt.Printf("Can't receive: %v\n", err)
			break
		}
		fmt.Println(reply.Text)
		if err := websocket.JSON.Send(ws, reply); err != nil {
			fmt.Printf("Can't send: %v\n", err)
			break
		}

	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./websocket")))

	mux.Handle("/text", websocket.Handler(EchoText))
	mux.Handle("/json", websocket.Handler(EchoJSON))

	srv := &http.Server{
		Addr:           "localhost:" + "3000",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

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

	fmt.Println("Websocket App started!")
	log.Printf("Starting HTTP server on %s", "3000")

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped
	fmt.Println("Websocket App stopped")
}
