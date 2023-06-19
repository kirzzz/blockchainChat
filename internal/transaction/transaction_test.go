package transaction

import (
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha256"
	"crypto/x509"
	_ "encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	// Создание тестовых входных данных
	inputs := []MessageInput{
		{
			TransactionID: "transaction_id_1",
			OutputIndex:   0,
			EncryptedData: []byte("encrypted_data_1"),
		},
		{
			TransactionID: "transaction_id_2",
			OutputIndex:   1,
			EncryptedData: []byte("encrypted_data_2"),
		},
	}

	outputs := []MessageOutput{
		{
			EncryptedData: []byte("encrypted_data_3"),
			Recipient:     "recipient_1",
		},
		{
			EncryptedData: []byte("encrypted_data_4"),
			Recipient:     "recipient_2",
		},
	}

	// Создание новой транзакции
	tx, err := NewTransaction(inputs, outputs)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.NotEmpty(t, tx.ID)
	assert.Len(t, tx.Inputs, len(inputs))
	assert.Len(t, tx.Outputs, len(outputs))
}

func TestTransaction_EncryptMessages(t *testing.T) {
	// Создание тестовых входных данных
	inputs := []MessageInput{
		{
			TransactionID: "transaction_id_1",
			OutputIndex:   0,
			EncryptedData: []byte("unencrypted_data_1"),
		},
	}
	outputs := []MessageOutput{
		{
			EncryptedData: []byte("unencrypted_data_2"),
			Recipient:     "recipient_1",
		},
	}

	// Создание новой транзакции
	tx := &Transaction{
		ID:      "transaction_id",
		Inputs:  inputs,
		Outputs: outputs,
	}

	// Шифрование сообщений
	err := tx.EncryptMessages()
	assert.NoError(t, err)

	// Проверка, что данные сообщений были зашифрованы
	for _, output := range tx.Outputs {
		assert.NotEqual(t, []byte("unencrypted_data_2"), output.EncryptedData)
	}

	// Проверка, что данные сообщений не изменились для входов
	for _, input := range tx.Inputs {
		assert.Equal(t, []byte("unencrypted_data_1"), input.EncryptedData)
	}
}

func TestTransaction_DecryptMessages(t *testing.T) {
	// Создание тестовых входных данных
	inputs := []MessageInput{
		{
			TransactionID: "transaction_id_1",
			OutputIndex:   0,
			EncryptedData: []byte("encrypted_data_1"),
		},
	}
	outputs := []MessageOutput{
		{
			EncryptedData: []byte("encrypted_data_2"),
			Recipient:     "recipient_1",
		},
	}

	// Создание новой транзакции
	tx := &Transaction{
		ID:      "transaction_id",
		Inputs:  inputs,
		Outputs: outputs,
	}

	// Генерация и загрузка приватного ключа получателя
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	// Расшифровка сообщений
	err = tx.DecryptMessages(privateKey)
	assert.NoError(t, err)

	// Проверка, что данные сообщений были расшифрованы
	for _, input := range tx.Inputs {
		assert.NotEqual(t, []byte("encrypted_data_1"), input.EncryptedData)
	}

	// Проверка, что данные сообщений не изменились для выходов
	for _, output := range tx.Outputs {
		assert.Equal(t, []byte("encrypted_data_2"), output.EncryptedData)
	}
}

func TestSerializeDeserializeTransaction(t *testing.T) {
	// Создание тестовых входных данных
	inputs := []MessageInput{
		{
			TransactionID: "transaction_id_1",
			OutputIndex:   0,
			EncryptedData: []byte("encrypted_data_1"),
		},
	}
	outputs := []MessageOutput{
		{
			EncryptedData: []byte("encrypted_data_2"),
			Recipient:     "recipient_1",
		},
	}

	// Создание новой транзакции
	tx1 := &Transaction{
		ID:      "transaction_id",
		Inputs:  inputs,
		Outputs: outputs,
	}

	// Сериализация транзакции
	data, err := tx1.Serialize()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Десериализация транзакции
	tx2, err := DeserializeTransaction(data)
	assert.NoError(t, err)
	assert.NotNil(t, tx2)
	assert.Equal(t, tx1.ID, tx2.ID)
	assert.Len(t, tx2.Inputs, len(tx1.Inputs))
	assert.Len(t, tx2.Outputs, len(tx1.Outputs))
}

func TestLoadPublicKey(t *testing.T) {
	// Генерация и сохранение публичного ключа получателя в PEM-кодированном файле
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	publicKeyFile := "public.pem"
	err = SavePublicKey(publicKeyFile, &privateKey.PublicKey)
	assert.NoError(t, err)

	// Загрузка публичного ключа из файла
	publicKey, err := LoadPublicKey(publicKeyFile)
	assert.NoError(t, err)
	assert.NotNil(t, publicKey)

	// Проверка, что загруженный публичный ключ соответствует ожидаемому
	assert.Equal(t, privateKey.PublicKey.N, publicKey.N)
	assert.Equal(t, privateKey.PublicKey.E, publicKey.E)
}

func SavePublicKey(filename string, publicKey *rsa.PublicKey) error {
	publicKeyData, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyData,
	}

	pemData := pem.EncodeToMemory(pemBlock)
	if pemData == nil {
		return errors.New("failed to encode PEM data")
	}

	err = ioutil.WriteFile(filename, pemData, 0644)
	if err != nil {
		return err
	}

	return nil
}
