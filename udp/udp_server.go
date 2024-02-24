package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var (
	address = ":8080" // address by default - :8080
)

func main() {
	arguments := os.Args
	if len(arguments) == 2 {
		address = arguments[1]
	}

	s, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Инициализируем UDP слушателя
	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Server is running.")

	defer func(connection *net.UDPConn) {
		err := connection.Close()
		if err != nil {
			fmt.Printf("failed to listen on address %s: %s\n", address, err)
			return
		}
	}(connection)
	buffer := make([]byte, 1024)

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		fmt.Print("-> ", string(buffer[0:n-1]))

		_, err = connection.WriteToUDP(
			[]byte(
				fmt.Sprintf("Message was received at UDP server at %s\n", time.Now().Format(time.DateTime)),
			),
			addr,
		)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
