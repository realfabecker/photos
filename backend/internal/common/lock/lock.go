package lock

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func Base64Encode(value []byte) string {
	return base64.URLEncoding.EncodeToString(value)
}

func Base64Decode(value string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(value)
}

func Encrypt(value []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(fmt.Sprintf("%x", md5.Sum([]byte("wallet_"+key)))))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	crp := gcm.Seal(nonce, nonce, value, nil)
	return crp, nil
}

func Decrypt(ciphered []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(fmt.Sprintf("%x", md5.Sum([]byte("wallet_"+key)))))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nSize := gcm.NonceSize()
	nonce, cipheredText := ciphered[:nSize], ciphered[nSize:]
	return gcm.Open(nil, nonce, cipheredText, nil)
}
