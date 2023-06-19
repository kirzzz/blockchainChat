package blockchain

import (
	"encoding/json"
	"time"
)

type Block struct {
	Index        int64
	Timestamp    int64
	Data         string
	PrevHash     string
	Nonce        int64
	Hash         string
	Difficulty   int
	MinerAddress string
}

type Blockchain struct {
	Difficulty int
	Tip        []byte
	db         DbInterface
}

func NewBlock(index int64, timestamp int64, data string, prevHash string, difficulty int, minerAddress string) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    timestamp,
		Data:         data,
		PrevHash:     prevHash,
		Difficulty:   difficulty,
		MinerAddress: minerAddress,
	}

	pow := NewProofOfWork(block, difficulty)
	block.Hash, block.Nonce = pow.Run()
	return block
}

func NewBlockchain(difficulty int, dbStorage DbInterface) (*Blockchain, error) {
	blockchain := &Blockchain{
		Difficulty: difficulty,
		db:         dbStorage,
	}

	if dbStorage.BlockchainExistsInDB() {
		err := dbStorage.SetTipFromDB(blockchain)
		if err != nil {
			return nil, err
		}

		return blockchain, nil
	}

	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now().UnixNano(),
		Data:         "Genesis Block",
		PrevHash:     "",
		Difficulty:   difficulty,
		MinerAddress: "",
	}

	pow := NewProofOfWork(genesisBlock, difficulty)
	genesisBlock.Hash, genesisBlock.Nonce = pow.Run()

	err := dbStorage.SaveBlockToDB(genesisBlock)
	if err != nil {
		return nil, err
	}

	err = dbStorage.SaveTipToDB(genesisBlock.Hash)
	if err != nil {
		return nil, err
	}

	blockchain.Tip = []byte(genesisBlock.Hash)

	return blockchain, nil
}

func (bc *Blockchain) AddBlock(data string, minerAddress string) error {
	prevBlockData, err := bc.db.GetBlockFromDB(string(bc.Tip))
	if err != nil {
		return err
	}

	var prevBlock Block
	err = json.Unmarshal(prevBlockData, &prevBlock)
	if err != nil {
		return err
	}

	newBlock := NewBlock(prevBlock.Index+1, time.Now().UnixNano(), data, prevBlock.Hash, bc.Difficulty, minerAddress)
	err = bc.db.SaveBlockToDB(newBlock)
	if err != nil {
		return err
	}

	err = bc.db.SaveTipToDB(newBlock.Hash)
	if err != nil {
		return err
	}

	bc.Tip = []byte(newBlock.Hash)

	return nil
}

// Serialize сериализует блок в байтовый массив
func (block *Block) Serialize() ([]byte, error) {
	data := make(map[string]interface{})
	data["Index"] = block.Index
	data["Timestamp"] = block.Timestamp
	data["Data"] = block.Data
	data["PrevHash"] = block.PrevHash
	data["Nonce"] = block.Nonce
	data["Hash"] = block.Hash
	data["Difficulty"] = block.Difficulty
	data["MinerAddress"] = block.MinerAddress

	return json.Marshal(data)
}
