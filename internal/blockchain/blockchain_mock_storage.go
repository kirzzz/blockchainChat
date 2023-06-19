package blockchain

import (
	"fmt"
)

// MockDbStorage Создаем мокированное хранилище данных для тестирования
type MockDbStorage struct {
	blockchainExistsInDB bool
	blockchainTip        []byte
	blocks               map[string][]byte
}

func NewMockDbStorage() *MockDbStorage {
	return &MockDbStorage{
		blockchainExistsInDB: false,
		blockchainTip:        nil,
		blocks:               make(map[string][]byte),
	}
}

func (db *MockDbStorage) BlockchainExistsInDB() bool {
	return db.blockchainExistsInDB
}

func (db *MockDbStorage) SetTipFromDB(blockchain *Blockchain) error {
	if db.blockchainTip == nil {
		return nil
	}

	blockchain.Tip = db.blockchainTip
	return nil
}

func (db *MockDbStorage) SaveBlockToDB(block *Block) error {
	serialized, err := block.Serialize()
	if err != nil {
		return err
	}

	db.blocks[block.Hash] = serialized
	return nil
}

func (db *MockDbStorage) SaveTipToDB(tip string) error {
	db.blockchainTip = []byte(tip)
	return nil
}

func (db *MockDbStorage) GetBlockFromDB(hash string) ([]byte, error) {
	blockData, exists := db.blocks[hash]
	if !exists {
		return nil, fmt.Errorf("block not found in DB")
	}

	return blockData, nil
}
