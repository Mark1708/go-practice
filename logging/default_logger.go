package main

import (
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("logging/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Настраиваем вывод в файл
	log.SetOutput(file)
	// Настраиваем формат логирования
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC | log.Lshortfile)

	log.Print("Logging to a file in Go!")
	// 2024/03/02 14:42:41.300185 custom_writer.go:18: Logging to a file in Go!
}
