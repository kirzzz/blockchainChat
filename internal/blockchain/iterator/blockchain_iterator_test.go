package iterator

import (
	"blockchainStorage/internal/blockchain"
	"testing"
)

func TestBlockchainIterator_Next(t *testing.T) {
	// Создаем фейковую базу данных для тестов
	db := blockchain.NewMockDbStorage()

	// Создаем несколько блоков для тестов
	block1 := &blockchain.Block{
		Index:     1,
		Timestamp: 123456789,
		Data:      "Block 1",
		PrevHash:  "prevhash1",
		Hash:      "hash1",
	}
	block2 := &blockchain.Block{
		Index:     2,
		Timestamp: 123456790,
		Data:      "Block 2",
		PrevHash:  "hash1",
		Hash:      "hash2",
	}
	block3 := &blockchain.Block{
		Index:     3,
		Timestamp: 123456791,
		Data:      "Block 3",
		PrevHash:  "hash2",
		Hash:      "hash3",
	}

	db.SaveBlockToDB(block1)
	db.SaveBlockToDB(block2)
	db.SaveBlockToDB(block3)

	// Создаем BlockchainIterator с текущим хэшем последнего блока
	iterator := NewBlockchainIterator(db, []byte(block3.Hash))

	// Проверяем, что следующий блок возвращает блок3, затем блок2, и затем блок1
	// в соответствии с порядком цепочки блоков
	validateBlock(t, iterator, block3)
	validateBlock(t, iterator, block2)
	validateBlock(t, iterator, block1)

	// Проверяем, что после достижения генезис-блока (первого блока) Next() возвращает ошибку
	_, err := iterator.Next()
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}

func validateBlock(t *testing.T, iterator *BlockchainIterator, expectedBlock *blockchain.Block) {
	block, err := iterator.Next()
	if err != nil {
		t.Errorf("Next() returned an error: %s, expect hash %s", err, string(iterator.currentHash))
		return
	}

	if block.Index != expectedBlock.Index ||
		block.Timestamp != expectedBlock.Timestamp ||
		block.Data != expectedBlock.Data ||
		block.PrevHash != expectedBlock.PrevHash ||
		block.Hash != expectedBlock.Hash {
		t.Errorf("Expected block %+v, but got %+v", expectedBlock, block)
	}
}
