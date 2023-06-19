package storage

import (
	"blockchainStorage/internal/blockchain"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

var dbMutex sync.Mutex

func setupDataStore(t *testing.T) (*DataStore, func()) {
	dbPath := "./test_db"

	// Создаем временную папку для БД LevelDB
	err := os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create DB directory: %v", err)
	}

	// Инициализируем хранилище данных
	ds, err := NewDataStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create DataStore: %v", err)
	}

	// Функция для очистки временной БД LevelDB
	cleanupDB := func() {
		err := ds.Close()
		if err != nil {
			t.Fatalf("failed to close DataStore: %v", err)
		}

		err = os.RemoveAll(dbPath)
		if err != nil {
			t.Fatalf("failed to remove DB directory: %v", err)
		}
	}

	return ds, cleanupDB
}

func TestPutAndGet(t *testing.T) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Инициализируем хранилище данных
	ds, cleanupDB := setupDataStore(t)
	defer cleanupDB()

	// Значение для сохранения
	value := "test value"

	// Сохраняем значение
	err := ds.Put("key", value)
	if err != nil {
		t.Fatalf("failed to put data: %v", err)
	}

	// Получаем значение по ключу
	var retrievedValue string
	_, err = ds.Get("key", &retrievedValue)
	if err != nil {
		t.Fatalf("failed to get data: %v", err)
	}

	// Проверяем, что полученное значение совпадает с сохраненным
	if retrievedValue != value {
		t.Errorf("retrieved value doesn't match the original value: got %s, expected %s", retrievedValue, value)
	}
}

func TestBlockchainExistsInDB(t *testing.T) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Инициализируем хранилище данных
	ds, cleanupDB := setupDataStore(t)
	defer cleanupDB()

	// Создаем блокчейн и сохраняем флаг существования в БД
	err := ds.Put(BlockchainExistsKey, true)
	if err != nil {
		t.Fatalf("failed to put blockchain existence flag: %v", err)
	}

	// Проверяем, что блокчейн существует в БД
	exists := ds.BlockchainExistsInDB()
	if !exists {
		t.Error("expected blockchain existence in DB")
	}
}

func TestSaveAndGetBlock(t *testing.T) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Инициализируем хранилище данных
	ds, cleanupDB := setupDataStore(t)
	defer cleanupDB()

	// Создаем блок для сохранения
	block := &blockchain.Block{
		Index:      1,
		Timestamp:  time.Now().UnixNano(),
		Data:       "test data",
		PrevHash:   "",
		Hash:       "block123",
		Difficulty: 1,
	}

	// Сохраняем блок в БД
	err := ds.SaveBlockToDB(block)
	if err != nil {
		t.Fatalf("failed to save block to DB: %v", err)
	}

	// Получаем блок из БД по его хэшу
	var blockData blockchain.Block
	blockData, err = ds.GetBlockFromDB(block.Hash)
	if err != nil {
		t.Fatalf("failed to get block from DB: %v", err)
	}

	// Проверяем, что полученные данные блока совпадают с ожидаемыми значениями
	expectedBlock := &blockchain.Block{
		Index:      1,
		Timestamp:  block.Timestamp,
		Data:       "test data",
		PrevHash:   "",
		Hash:       "block123",
		Difficulty: 1,
	}

	if !reflect.DeepEqual(&blockData, expectedBlock) {
		t.Errorf("retrieved block data doesn't match the expected data: got %+v, expected %+v", &blockData, expectedBlock)
	}
}

func TestSaveAndRetrieveTip(t *testing.T) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Инициализируем хранилище данных
	ds, cleanupDB := setupDataStore(t)
	defer cleanupDB()

	// Хэш последнего блока
	tipHash := "block123"

	// Сохраняем хэш последнего блока в БД
	err := ds.SaveTipToDB(tipHash)
	if err != nil {
		t.Fatalf("failed to save tip to DB: %v", err)
	}

	// Получаем хэш последнего блока из БД
	var retrievedTipHash string
	_, err = ds.Get(TipKey, &retrievedTipHash)
	if err != nil {
		t.Fatalf("failed to get tip from DB: %v", err)
	}

	// Проверяем, что полученный хэш совпадает с ожидаемым
	if retrievedTipHash != tipHash {
		t.Errorf("retrieved tip hash doesn't match the expected hash: got %s, expected %s", retrievedTipHash, tipHash)
	}
}
