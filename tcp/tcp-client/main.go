package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Запущен TCP клиент")

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Пожалуйта укажите хост и порт в формате: host:port.")
		return
	}

	// Подключаемся к TCP серверу
	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

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
			fmt.Println(err)
			return
		}
		fmt.Print("-> " + message)
	}
}
