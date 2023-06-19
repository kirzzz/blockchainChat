package blockchain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewBlockchain(t *testing.T) {
	dbStorage := NewMockDbStorage() // Создаем экземпляр мокированного хранилища данных
	difficulty := 3

	bc, err := NewBlockchain(difficulty, dbStorage)
	if err != nil {
		t.Fatalf("failed to create new blockchain: %v", err)
	}

	// Проверяем, что сложность установлена корректно
	if bc.Difficulty != difficulty {
		t.Errorf("expected difficulty %d, but got %d", difficulty, bc.Difficulty)
	}

	// Проверяем, что создан блок Genesis и сохранен в БД
	genesisBlockData, err := dbStorage.GetBlockFromDB(string(bc.Tip))
	if err != nil {
		t.Fatalf("failed to get genesis block from DB: %v", err)
	}

	var genesisBlock Block
	err = json.Unmarshal(genesisBlockData, &genesisBlock)
	if err != nil {
		t.Fatalf("failed to unmarshal genesis block data: %v", err)
	}

	// Проверяем, что индекс Genesis блока равен 0
	if genesisBlock.Index != 0 {
		t.Errorf("expected genesis block index 0, but got %d", genesisBlock.Index)
	}

	// Проверяем, что предыдущий хэш Genesis блока пустой
	if genesisBlock.PrevHash != "" {
		t.Errorf("expected empty prevHash for genesis block, but got %s", genesisBlock.PrevHash)
	}
}

func TestAddBlock(t *testing.T) {
	dbStorage := NewMockDbStorage() // Создаем экземпляр мокированного хранилища данных
	difficulty := 3
	minerAddress := "miner_address"

	bc, err := NewBlockchain(difficulty, dbStorage)
	if err != nil {
		t.Fatalf("failed to create new blockchain: %v", err)
	}

	// Добавляем блок в блокчейн
	data := "Block Data"
	err = bc.AddBlock(data, minerAddress)
	if err != nil {
		t.Fatalf("failed to add block to blockchain: %v", err)
	}

	// Проверяем, что блок добавлен в БД
	blockData, err := dbStorage.GetBlockFromDB(string(bc.Tip))
	if err != nil {
		t.Fatalf("failed to get block from DB: %v", err)
	}

	var block Block
	err = json.Unmarshal(blockData, &block)
	if err != nil {
		t.Fatalf("failed to unmarshal block data: %v", err)
	}

	// Проверяем, что индекс нового блока увеличился на 1
	expectedIndex := int64(1)
	if block.Index != expectedIndex {
		t.Errorf("expected block index %d, but got %d", expectedIndex, block.Index)
	}

	// Проверяем, что предыдущий хэш нового блока соответствует хэшу последнего блока в блокчейне
	lastBlockData, err := dbStorage.GetBlockFromDB(block.PrevHash)
	if err != nil {
		t.Fatalf("failed to get last block from DB: %v", err)
	}

	var lastBlock Block
	err = json.Unmarshal(lastBlockData, &lastBlock)
	if err != nil {
		t.Fatalf("failed to unmarshal last block data: %v", err)
	}

	if block.PrevHash != lastBlock.Hash {
		t.Errorf("expected previous hash %s, but got %s", lastBlock.Hash, block.PrevHash)
	}
}

func TestSerialize(t *testing.T) {
	block := &Block{
		Index:        1,
		Timestamp:    time.Now().UnixNano(),
		Data:         "Block Data",
		PrevHash:     "Previous Hash",
		Nonce:        12345,
		Hash:         "Block Hash",
		Difficulty:   3,
		MinerAddress: "miner_address",
	}

	serialized, err := block.Serialize()
	if err != nil {
		t.Fatalf("failed to serialize block: %v", err)
	}

	// Проверяем, что сериализованные данные содержат необходимые поля
	var data map[string]interface{}
	err = json.Unmarshal(serialized, &data)
	if err != nil {
		t.Fatalf("failed to unmarshal serialized data: %v", err)
	}

	expectedFields := []string{"Index", "Timestamp", "Data", "PrevHash", "Nonce", "Hash", "Difficulty", "MinerAddress"}
	for _, field := range expectedFields {
		_, exists := data[field]
		if !exists {
			t.Errorf("expected field %s in serialized data, but it is missing", field)
		}
	}
}
