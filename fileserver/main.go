package main

import (
	"log"
	"net/http"
)

func main() {
	// Передача содержимого директории ./frontend по адресу Route /.
	http.Handle("/", http.FileServer(http.Dir("./fileserver/frontend")))

	// ЗАпуск сервера
	log.Fatal(http.ListenAndServe(":8080", nil))
}
