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
		panic("pass in args host and port of UDP server - host:port.")
		return
	}

	address := arguments[1]

	// Устанавливаем адрес для прослушивания
	s, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Подключаемся к UDP серверу
	c, err := connect(s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text + "\n")
		_, err = c.Write(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Читаем ответ сервера и выводим
		buffer := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			conn, err := connect(s)
			if err != nil {
				return
			}
			c = conn
		} else {
			fmt.Print("-> " + string(buffer[0:n]))
		}
	}
}

func connect(s *net.UDPAddr) (*net.UDPConn, error) {
	var lastError error
	for retryCount := 0; retryCount < mrc; retryCount++ {
		c, err := net.DialUDP("udp4", nil, s)
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
