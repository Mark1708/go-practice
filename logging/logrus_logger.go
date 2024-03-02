package main

import (
	"github.com/google/uuid"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	file, err := os.OpenFile("logging/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{}) // Также есть log.TextFormatter{}
	log.SetLevel(log.WarnLevel)

	log.WithFields(log.Fields{
		"spanId":    uuid.New().String(),
		"requestId": uuid.New().String(),
	}).Warn("Some Logging Message 1")
	// {"level":"warning","msg":"Some Logging Message 1","requestId":"bcb34cc2-fe1d-4d9b-bc4a-2d770ba5e825","spanId":"dad3742d-85a1-49a8-97e3-d3ff4361b122","time":"2024-03-02T18:24:38+03:00"}

	log.Warn("Some Logging Message 2")
	// {"level":"warning","msg":"Some Logging Message 2","time":"2024-03-02T18:24:38+03:00"}
}
