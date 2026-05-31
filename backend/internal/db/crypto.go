package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

func getEncryptionKey() []byte {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		keyStr = "majhigov-default-super-secret-key-2026"
	}
	hash := sha256.Sum256([]byte(keyStr))
	return hash[:]
}

// Encrypt encrypts plain text string to hex-encoded string using AES-GCM
func Encrypt(plainText string) (string, error) {
	if plainText == "" {
		return "", nil
	}
	key := getEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts hex-encoded string to plain text string using AES-GCM
func Decrypt(cipherTextHex string) (string, error) {
	if cipherTextHex == "" {
		return "", nil
	}
	ciphertext, err := hex.DecodeString(cipherTextHex)
	if err != nil {
		return "", err
	}
	key := getEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
