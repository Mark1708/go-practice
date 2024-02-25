package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopractice/mqtt/helper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Настраиваем клиента
	opts := helper.GetBaseClientOpt()
	opts.SetDefaultPublishHandler(defaultHandler)
	client := mqtt.NewClient(opts)

	// Подключаемся
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Subscriber Started")

	// Подписываемся на топики
	go sub(client, "example/stdin/#", nil)
	go sub(client, "example/simple", simpleHandler)
	go sub(client, "example/multi/#", multiHandler)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Disconnect(250)

	fmt.Println("Subscriber Disconnected")
}

// sub Функция подписки на топик
func sub(client mqtt.Client, topic string, callback mqtt.MessageHandler) {
	if token := client.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

// defaultHandler Функция обработчик сообщений для не обрабатываемых топиков (дефолтный обработчик)
// На эти топики всё равно требуется подписка
func defaultHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("DefaultHandler   ")
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}

// simpleHandler Функция обработчик сообщений для топика example/simple
func simpleHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("SimpleHandler   ")
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}

// multiHandler Функция обработчик сообщений для топиков example/multi/#
func multiHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("MultiHandler   ")
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}
