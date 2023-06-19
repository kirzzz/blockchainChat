package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

type ProofOfWork struct {
	block      *Block
	difficulty int
}

func NewProofOfWork(block *Block, difficulty int) *ProofOfWork {
	return &ProofOfWork{
		block:      block,
		difficulty: difficulty,
	}
}

func (pow *ProofOfWork) Run() (string, int64) {
	target := bytes.Repeat([]byte("0"), pow.difficulty)
	var hash string
	var nonce int64 = 0

	for !isValidHash(hash, string(target)) {
		hash = pow.calculateHash(nonce)
		nonce++
	}

	return hash, nonce - 1
}

func isValidHash(hash string, target string) bool {
	return len(hash) >= len(target) &&
		hash[:len(target)] == target
}

func (pow *ProofOfWork) calculateHash(nonce int64) string {
	record := strconv.FormatInt(pow.block.Index, 10) +
		strconv.FormatInt(pow.block.Timestamp, 10) +
		pow.block.Data +
		pow.block.PrevHash +
		strconv.FormatInt(nonce, 10) +
		pow.block.MinerAddress

	h := sha256.New()
	h.Write([]byte(record))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
