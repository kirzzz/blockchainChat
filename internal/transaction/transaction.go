package transaction

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
)

type MessageInput struct {
	TransactionID string
	OutputIndex   int
	EncryptedData []byte
}

type MessageOutput struct {
	EncryptedData []byte
	Recipient     string
}

type Transaction struct {
	ID      string
	Inputs  []MessageInput
	Outputs []MessageOutput
}

// NewTransaction создает новую транзакцию с зашифрованными сообщениями
func NewTransaction(inputs []MessageInput, outputs []MessageOutput) (*Transaction, error) {
	tx := &Transaction{
		ID:      generateTransactionID(),
		Inputs:  inputs,
		Outputs: outputs,
	}

	err := tx.EncryptMessages()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// generateTransactionID генерирует уникальный идентификатор транзакции
func generateTransactionID() string {
	id := uuid.New()
	return id.String()
}

// EncryptMessages шифрует сообщения в транзакции с помощью публичных ключей получателей
func (tx *Transaction) EncryptMessages() error {
	for i := range tx.Outputs {
		recipientPublicKey, err := LoadPublicKey(tx.Outputs[i].Recipient)
		if err != nil {
			return fmt.Errorf("failed to load public key for recipient %s: %w", tx.Outputs[i].Recipient, err)
		}

		encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, recipientPublicKey, tx.Outputs[i].EncryptedData, nil)
		if err != nil {
			return fmt.Errorf("failed to encrypt message data for recipient %s: %w", tx.Outputs[i].Recipient, err)
		}

		tx.Outputs[i].EncryptedData = encryptedData
	}

	return nil
}

// DecryptMessages расшифровывает сообщения в транзакции с помощью приватного ключа получателя
func (tx *Transaction) DecryptMessages(privateKey *rsa.PrivateKey) error {
	for i := range tx.Inputs {
		decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, tx.Inputs[i].EncryptedData, nil)
		if err != nil {
			return fmt.Errorf("failed to decrypt message data for transaction %s: %w", tx.ID, err)
		}

		tx.Inputs[i].EncryptedData = decryptedData
	}

	return nil
}

// Serialize сериализует транзакцию в байтовый массив
func (tx *Transaction) Serialize() ([]byte, error) {
	return json.Marshal(tx)
}

// DeserializeTransaction десериализует транзакцию из байтового массива
func DeserializeTransaction(data []byte) (*Transaction, error) {
	var tx Transaction
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// LoadPublicKey загружает публичный ключ получателя из PEM-кодированного файла
func LoadPublicKey(publicKeyFile string) (*rsa.PublicKey, error) {
	publicKeyData, err := ReadPEMFile(publicKeyFile)
	if err != nil {
		return nil, err
	}

	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return publicKey, nil
}

// ReadPEMFile читает PEM-кодированные данные из файла
func ReadPEMFile(filename string) ([]byte, error) {
	pemData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read PEM file: %w", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to decode PEM file")
	}

	return block.Bytes, nil
}
