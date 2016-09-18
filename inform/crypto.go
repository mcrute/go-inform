package inform

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func Pad(src []byte, blockSize int) []byte {
	padLen := blockSize - (len(src) % blockSize)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padText...)
}

func Unpad(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src)
	paddingLen := int(src[srcLen-1])
	if paddingLen >= srcLen || paddingLen > blockSize {
		return nil, errors.New("Padding size error")
	}
	return src[:srcLen-paddingLen], nil
}

func decodeHexKey(key string) (cipher.Block, error) {
	decodedKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func makeAESIV() ([]byte, error) {
	iv := make([]byte, 16)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	return iv, nil
}

// Returns ciphertext and IV, does not modify payload
func Encrypt(payload []byte, key string) ([]byte, []byte, error) {
	ct := make([]byte, len(payload))
	copy(ct, payload)
	ct = Pad(ct, aes.BlockSize)

	iv, err := makeAESIV()
	if err != nil {
		return nil, nil, err
	}

	block, err := decodeHexKey(key)
	if err != nil {
		return nil, nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ct, ct)

	return ct, iv, nil
}

func Decrypt(payload, iv []byte, key string) ([]byte, error) {
	b := make([]byte, len(payload))

	block, err := decodeHexKey(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(b, payload)

	u, err := Unpad(b, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	return u, nil
}
