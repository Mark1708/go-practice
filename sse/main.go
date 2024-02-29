package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./sse")))

	http.HandleFunc("/time", timeHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	/*
	 * Сервер должен ответить со статусом 200 и заголовком Content-Type: text/event-stream,
	 * затем он должен поддерживать соединение открытым и отправлять сообщения в особом формате
	 */
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	counter := 0
	for {
		now := time.Now().Format(time.TimeOnly)
		counter++
		/*
		 * Текст сообщения указывается после data:, пробел после двоеточия необязателен.
		 * Сообщения разделяются двойным переносом строки \n\n.
		 * Чтобы разделить сообщение на несколько строк, мы можем отправить несколько data: подряд (третье сообщение).
		 * Сервер может выставить рекомендуемую задержку, указав в ответе retry: (в миллисекундах).
		 * Чтобы правильно возобновить подключение, каждое сообщение должно иметь поле id.
		 *
		 * Типы событий
		 * По умолчанию объект EventSource генерирует 3 события:
		 * 	message – получено сообщение, доступно как event.data.
		 * 	open – соединение открыто.
		 * 	error – не удалось установить соединение, например, сервер вернул статус 500.
		 * Сервер может указать другой тип события с помощью event: ... в начале сообщения.
		 */
		fmt.Fprintf(w, "id: %d\nevent: %s\nretry: %d\ndata: %s\n\n", counter, "test", 5*time.Millisecond, now)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(time.Second)
	}
}
