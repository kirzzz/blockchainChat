package network

import (
	"blockchainStorage/common"
	"encoding/json"
	"net"
	"testing"
)

func TestBroadcast(t *testing.T) {
	// Создаем виртуальный сервер
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	// Запускаем сервер в отдельной горутине
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Fatalf("Failed to accept connection: %v", err)
		}
		defer conn.Close()

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatalf("Failed to read data: %v", err)
		}

		var receivedMsg Message
		err = json.Unmarshal(buf[:n], &receivedMsg)
		if err != nil {
			t.Fatalf("Failed to unmarshal message: %v", err)
		}

		expectedMsg := Message{
			Command: "testCommand",
			Data:    []byte("testData"),
		}
		if receivedMsg.Command != expectedMsg.Command || string(receivedMsg.Data) != string(expectedMsg.Data) {
			t.Errorf("Received message: got %+v, expected %+v", receivedMsg, expectedMsg)
		}
	}()

	// Подготавливаем список узлов
	nodeList := []common.Node{
		{Address: listener.Addr().String()},
	}

	// Создаем экземпляр Network
	network := &Network{
		NodeList: nodeList,
	}

	// Отправляем данные
	err = network.Broadcast("testCommand", []byte("testData"))
	if err != nil {
		t.Fatalf("Failed to broadcast data: %v", err)
	}
}

func TestStartServer(t *testing.T) {
	// Создаем виртуальный сервер
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	// Подготавливаем сообщение для отправки
	msgToSend := Message{
		Command: "testCommand",
		Data:    []byte("testData"),
	}
	payload, err := json.Marshal(msgToSend)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Fatalf("Failed to accept connection: %v", err)
		}
		defer conn.Close()

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatalf("Failed to read data: %v", err)
		}

		var receivedMsg Message
		err = json.Unmarshal(buf[:n], &receivedMsg)
		if err != nil {
			t.Fatalf("Failed to unmarshal message: %v", err)
		}

		expectedMsg := msgToSend
		if receivedMsg.Command != expectedMsg.Command || string(receivedMsg.Data) != string(expectedMsg.Data) {
			t.Errorf("Received message: got%+v, expected %+v", receivedMsg, expectedMsg)
		}
	}()
	// Создаем экземпляр Network
	network := &Network{}

	// Запускаем сервер
	go func() {
		err := network.StartServer(listener.Addr().(*net.TCPAddr).Port, func(msg *Message, conn net.Conn) {
			// Обработка сообщения
		})
		if err != nil {
			t.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Подключаемся к виртуальному серверу
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Отправляем данные
	_, err = conn.Write(payload)
	if err != nil {
		t.Fatalf("Failed to send data: %v", err)
	}
}
