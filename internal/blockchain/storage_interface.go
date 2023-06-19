package blockchain

type DbInterface interface {
	BlockchainExistsInDB() bool
	SetTipFromDB(blockchain *Blockchain) error
	SaveBlockToDB(block *Block) error
	SaveTipToDB(tipHash string) error
	GetBlockFromDB(blockHash string) ([]byte, error)
}
