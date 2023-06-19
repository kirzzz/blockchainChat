package key_gen

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

func GenerateKey() {
	// Выбор эллиптической кривой (например, P-256)
	curve := elliptic.P256()

	// Генерация приватного ключа
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println("Ошибка генерации приватного ключа:", err)
		return
	}

	// Получение публичного ключа из приватного
	publicKey := &privateKey.PublicKey

	// Вывод приватного и публичного ключей
	fmt.Printf("Приватный ключ: %x\n", privateKey.D)
	fmt.Printf("Публичный ключ: %x\n", elliptic.Marshal(curve, publicKey.X, publicKey.Y))
}
