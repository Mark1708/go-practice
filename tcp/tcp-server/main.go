package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	fmt.Println("Запущен TCP сервер")

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Пожалуйста укажите порт")
		return
	}

	PORT := ":" + arguments[1]
	// Инициализируем прослушивание порта
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Закрываем прослушивание после завершения работы
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			panic("error during closing listener")
		}
	}(l)

	var connMap = &sync.Map{}

	for {
		// Принимаем подключение
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		id := uuid.New().String()
		connMap.Store(id, c)
		go handleMessage(id, c, connMap)
	}
}

func handleMessage(id string, c net.Conn, connMap *sync.Map) {
	defer func(c net.Conn) {
		err := c.Close()
		connMap.Delete(id)
		if err != nil {
			panic("error during closing connection")
		}
	}(c)

	for {
		// Читаем сообщение
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		// Выводим сообщение
		fmt.Printf("[%s] -> %s", id, netData)

		// Отправляем время получения
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		_, _ = c.Write([]byte(myTime))

		// Рассылка клиентам
		//connMap.Range(func(key, value any) bool {
		//	if conn, ok := value.(net.Conn); ok {
		//
		//	}
		//})
	}
}
