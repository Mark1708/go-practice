package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// upgrader Структура с набором методов для работы с Websocket
// Например метод Upgrade для апгрейда подключения до уровня Websocket добавляет все необходимые заголовки
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// EchoText2 Метод получающий и отправляющий обычные текстовые сообщения
func EchoText2(w http.ResponseWriter, r *http.Request) {
	// обновление соединения до WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// цикл обработки сообщений
	for {
		messageType, message, err := ws.ReadMessage()

		if err != nil {
			fmt.Printf("Can't receive: %v\n", err)
			break
		}

		if err := ws.WriteMessage(messageType, message); err != nil {
			fmt.Printf("Can't send: %v\n", err)
			break
		}

	}
}

// EchoJSON2 Метод получающий и отправляющий JSON сообщения
func EchoJSON2(w http.ResponseWriter, r *http.Request) {
	// обновление соединения до WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// цикл обработки сообщений
	for {
		var reply Message

		if err := ws.ReadJSON(&reply); err != nil {
			fmt.Printf("Can't receive: %v\n", err)
			break
		}
		fmt.Println(reply.Text)
		if err := ws.WriteJSON(reply); err != nil {
			fmt.Printf("Can't send: %v\n", err)
			break
		}

	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./websocket")))

	mux.HandleFunc("/text", EchoText2)
	mux.HandleFunc("/json", EchoJSON2)

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
