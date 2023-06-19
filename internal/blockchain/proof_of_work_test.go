package blockchain

import (
	"strings"
	"testing"
)

func TestProofOfWork(t *testing.T) {
	block := &Block{
		Index:        1,
		Timestamp:    1234567800,
		Data:         "Test data",
		PrevHash:     "Prev Hash",
		Nonce:        0,
		Hash:         "",
		Difficulty:   5,
		MinerAddress: "Miner Address",
	}

	pow := NewProofOfWork(block, block.Difficulty)
	hash, nonce := pow.Run()

	// Verify that the hash has the required difficulty level
	target := strings.Repeat("0", pow.difficulty)
	if !isValidHash(hash, target) {
		t.Errorf("Proof of work failed. Expected hash with %d leading zeros, got: %s", pow.difficulty, hash)
	}

	// Verify that the nonce is correctly updated
	expectedNonce := int64(1915683) // Adjust the expected nonce value based on the specific difficulty level
	if nonce != expectedNonce {
		t.Errorf("Incorrect nonce. Expected: %d, got: %d", expectedNonce, nonce)
	}
}
