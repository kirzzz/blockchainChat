package iterator

import (
	"blockchainStorage/internal/blockchain"
	"encoding/json"
)

type BlockchainIterator struct {
	currentHash []byte
	db          blockchain.DbInterface
}

func NewBlockchainIterator(db blockchain.DbInterface, currentHash []byte) *BlockchainIterator {
	return &BlockchainIterator{
		currentHash: currentHash,
		db:          db,
	}
}

func (it *BlockchainIterator) Next() (*blockchain.Block, error) {
	blockData, err := it.db.GetBlockFromDB(string(it.currentHash))
	if err != nil {
		return nil, err
	}

	var block blockchain.Block
	err = json.Unmarshal(blockData, &block)
	if err != nil {
		return nil, err
	}

	it.currentHash = []byte(block.PrevHash)
	return &block, nil
}
