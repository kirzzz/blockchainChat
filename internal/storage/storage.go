package storage

import (
	"blockchainStorage/internal/blockchain"
	"encoding/json"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	BlockchainExistsKey = "blockchain_exists"
	TipKey              = "tip"
	BlockPrefix         = "block_"
)

// DataStore contracts/DbInterface

type DataStore struct {
	db            *leveldb.DB
	blockchainTip []byte
}

func (ds *DataStore) SetTipFromDB(blockchain *blockchain.Blockchain) error {
	if ds.blockchainTip == nil {
		return nil
	}

	blockchain.Tip = ds.blockchainTip
	return nil
}

func NewDataStore(dbPath string) (*DataStore, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open LevelDB: %w", err)
	}

	return &DataStore{
		db: db,
	}, nil
}

func (ds *DataStore) Close() error {
	err := ds.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close LevelDB: %w", err)
	}
	return nil
}

func (ds *DataStore) Put(key string, value interface{}) error {
	dataBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = ds.db.Put([]byte(key), dataBytes, nil)
	if err != nil {
		return fmt.Errorf("failed to put data in LevelDB: %w", err)
	}

	return nil
}

func (ds *DataStore) Get(key string, value interface{}) ([]byte, error) {
	dataBytes, err := ds.db.Get([]byte(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, fmt.Errorf("key not found")
		}
		return nil, fmt.Errorf("failed to get data from LevelDB: %w", err)
	}

	err = json.Unmarshal(dataBytes, &value)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return dataBytes, nil
}

func (ds *DataStore) Delete(key string) error {
	err := ds.db.Delete([]byte(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to delete data from LevelDB: %w", err)
	}

	return nil
}

// BlockchainExistsInDB Проверяет наличие блокчейна в БД
func (ds *DataStore) BlockchainExistsInDB() bool {
	var val interface{}
	_, err := ds.Get(BlockchainExistsKey, &val)
	return err == nil
}

// SaveBlockToDB Сохраняет блок в БД
func (ds *DataStore) SaveBlockToDB(block *blockchain.Block) error {
	err := ds.Put(BlockPrefix+block.Hash, block)
	if err != nil {
		return fmt.Errorf("failed to save block to DB: %w", err)
	}

	return nil
}

// SaveTipToDB Сохраняет хэш последнего блока в БД
func (ds *DataStore) SaveTipToDB(tipHash string) error {
	err := ds.Put(TipKey, tipHash)
	if err != nil {
		return fmt.Errorf("failed to save tip to DB: %w", err)
	}

	return nil
}

// GetBlockFromDB Получает блок из БД по его хэшу
func (ds *DataStore) GetBlockFromDB(blockHash string) ([]byte, error) {
	blockData, err := ds.Get(BlockPrefix+blockHash, nil)
	if err != nil {
		return []byte{}, err
	}

	return blockData, nil
}
