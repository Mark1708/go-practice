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

	// Устанавливаем адрес для прослушивания
	s, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Создаем соединение
	connection, err := net.ListenUDP("udp", s)
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

	// Буфер для хранения данных
	buffer := make([]byte, 1024)
	for {
		// Читаем сообщение
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Выводим сообщение
		fmt.Print("-> ", string(buffer[0:n-1]))

		// Отправляем время получения
		_, _ = connection.WriteToUDP(
			[]byte(
				fmt.Sprintf("Message was received at UDP server at %s\n", time.Now().Format(time.DateTime)),
			),
			addr,
		)
	}
}
