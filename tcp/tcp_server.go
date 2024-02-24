package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	address = ":8080" // address by default - :8080
)

// server Структура TCP сервера
type server struct {
	wg         sync.WaitGroup
	listener   net.Listener
	shutdown   chan struct{}
	connection chan net.Conn
	connMap    *sync.Map
}

// newServer Функция для создания сервера
func newServer(address string) (*server, error) {
	// Инициализируем TCP слушателя
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address %s: %w", address, err)
	}

	return &server{
		listener:   listener,
		shutdown:   make(chan struct{}),
		connection: make(chan net.Conn),
		connMap:    &sync.Map{},
	}, nil
}

// acceptConnections Процедура выполняемая в бесконечном цикле ждя прослушивания входящих соединений
func (s *server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		default:
			// Принимаем подключение
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}
			s.connection <- conn
		}
	}
}

// handleConnections Обработчик входящих соединений
func (s *server) handleConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		case conn := <-s.connection:
			go s.handleConnection(conn)
		}
	}
}

// handleConnection Обработчик соединения
func (s *server) handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)
	// Добавляем подключение
	id := uuid.New().String()
	s.connMap.Store(id, conn)

	for {
		// Читаем сообщение
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		// Выводим сообщение
		fmt.Printf("[%s] Received from %v -> %s", id, conn.RemoteAddr(), netData)

		// Отправляем время получения
		_, _ = conn.Write(
			[]byte(
				fmt.Sprintf("Message was received at TCP server at %s\n", time.Now().Format(time.DateTime)),
			),
		)
	}
}

// start Метод для запуска сервера
func (s *server) start() {
	s.wg.Add(2)
	go s.acceptConnections()
	go s.handleConnections()
}

// stop Метод для остановки сервера
func (s *server) stop() {
	close(s.shutdown)
	err := s.listener.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		fmt.Println("Timed out waiting for connections to finish")
		return
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 2 {
		address = arguments[1]
	}

	s, err := newServer(address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.start()
	fmt.Println("Server is running.")

	// Ждём сигнала SIGINT или SIGTERM, чтобы корректно завершить работу сервера
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down server...")
	s.stop()
	fmt.Println("Server stopped.")
}
