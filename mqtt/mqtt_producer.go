package main

import (
	"bufio"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopractice/mqtt/helper"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	// Настраиваем клиента
	opts := helper.GetBaseClientOpt()
	client := mqtt.NewClient(opts)

	// Подключаемся
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Producer Started")

	// Публикуем данные
	go pub(client, "example/simple", "Hello Simple Subscriber", 5)
	go pub(client, "example/multi/1", "Hello Multi Subscriber", 5)
	go pub(client, "example/multi/2", "Hello Multi Subscriber", 5)
	go stdinPub(client)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Disconnect(250)

	fmt.Println("Publisher Disconnected")
}

// stdinPub Функция позволяющая публиковать сообщения из консоли
// Для этого достаточно отправить: example/simple;Hello World! (первый аргумент - топик, а второй - сообщение)
func stdinPub(client mqtt.Client) {
	// Инициализируем reader из консоли
	reader := bufio.NewReader(os.Stdin)
	for {
		// Читаем сообщение из консоли
		fmt.Print(">> ")
		str, _ := reader.ReadString('\n')
		args := strings.Split(str, ";")
		token := client.Publish("example/stdin/"+args[0], 0, false, args[1])
		token.Wait()
	}
}

// pub Функция для публикации сообщений
func pub(client mqtt.Client, topic, payload string, num int) {
	for i := 0; i < num; i++ {
		token := client.Publish(topic, 0, false, payload)
		token.Wait()
		time.Sleep(time.Second)
	}
}
