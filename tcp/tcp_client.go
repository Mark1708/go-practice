package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	mrc = 3               // mrc Max reconnect count
	tbc = time.Second * 2 // tbc Time between connections
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		panic("pass in args host and port of TCP server - host:port.")
		return
	}

	// Подключаемся к TCP серверу
	address := arguments[1]
	c, err := connect(address)
	if err != nil {
		return
	}
	fmt.Printf("Client is running at address - %s.\n", c.LocalAddr())

	for {
		// Инициализируем reader из консоли
		reader := bufio.NewReader(os.Stdin)

		// Читаем сообщение из консоли
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		_, _ = fmt.Fprintf(c, text+"\n")

		// Читаем ответ сервера и выводим
		message, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			conn, err := connect(address)
			if err != nil {
				return
			}
			c = conn
		} else {
			fmt.Print("-> " + message)
		}
	}
}

func connect(address string) (net.Conn, error) {
	var lastError error
	for retryCount := 0; retryCount < mrc; retryCount++ {
		c, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Printf("Connection #%d broken with error: %v\n", retryCount+1, err)
			lastError = err
			time.Sleep(tbc)
		} else {
			return c, nil
		}
	}
	return nil, lastError
}
