package main

import (
	"blockchainStorage/config"
	"blockchainStorage/internal/blockchain"
	"blockchainStorage/internal/network"
	"blockchainStorage/internal/storage"
	"log"
	"net"
)

func main() {
	// Загрузка конфигурации из файла
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Создание и инициализация хранилища данных
	dataStore, err := storage.NewDataStore(cfg.DataBasePath)
	if err != nil {
		log.Fatal("Failed to initialize data store:", err)
	}
	defer dataStore.Close()

	// Создание и инициализация блокчейна
	chain, err := blockchain.NewBlockchain(cfg.Difficulty, dataStore)
	if err != nil {
		log.Fatal("Failed to initialize blockchain:", err)
	}

	// Создание и инициализация сети
	n := network.Network{NodeList: cfg.Nodes}

	// Запуск сервера для прослушивания входящих соединений
	go func() {
		err := n.StartServer(cfg.Port, handleIncomingMessage)
		if err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Отправка первичного блока Genesis всем узлам сети
	genesisBlock, _ := dataStore.GetBlockFromDB(string(chain.Tip))

	err = n.Broadcast("block", genesisBlock) //add serialize for block
	if err != nil {
		log.Println("Failed to broadcast genesis block:", err)
	}

	// Пример использования: сохранение данных в хранилище
	err = dataStore.Put("key", "value")
	if err != nil {
		log.Println("Failed to put data in data store:", err)
	}

	// Пример использования: отправка сообщения по сети
	message := "Hello, network!"
	err = n.Broadcast("message", []byte(message))
	if err != nil {
		log.Println("Failed to broadcast message:", err)
	}

	// Бесконечный цикл для работы приложения
	select {}
}

// Обработчик входящих сообщений
func handleIncomingMessage(msg *network.Message, conn net.Conn) {
	// Обработка входящего сообщения

	// Ваш код
}
