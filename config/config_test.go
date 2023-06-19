package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Создаем временный файл с тестовой конфигурацией
	tempFile, err := os.CreateTemp("", "test-config.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}(tempFile.Name())

	// Записываем тестовую конфигурацию в файл
	_, err = tempFile.WriteString(`{
		"port": 8080,
		"dbPath": "test.db",
		"nodes": ["peer1.example.com", "peer2.example.com"],
		"difficulty": 10,
		"blockReward": 5,
		"genesisBlockNonce": 20
	}`)

	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	err = tempFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Загружаем конфигурацию из временного файла
	config, err := LoadConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	expectedPort := 8080
	if config.Port != expectedPort {
		t.Errorf("Network.Port: got %d, expected %d", config.Port, expectedPort)
	}

	expectedPeers := []string{"peer1.example.com", "peer2.example.com"}
	for i, peer := range config.Nodes {
		if peer.Address != expectedPeers[i] {
			t.Errorf("Network.Peers[%d]: got %s, expected %s", i, peer, expectedPeers[i])
		}
	}
}
