package inform

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
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

func Decrypt(payload, iv []byte, key string, w *InformWrapper) ([]byte, error) {
	if !w.IsEncrypted() {
		return nil, errors.New("payload is not encrypted")
	}

	if w.IsGCMEncrypted() {
		return decryptGCM(payload, iv, key, w)
	} else {
		return decryptCBC(payload, iv, key)
	}

	return nil, nil
}

func buildAuthData(w *InformWrapper, iv []byte) []byte {
	ad := &bytes.Buffer{}
	binary.Write(ad, binary.BigEndian, int32(PROTOCOL_MAGIC))
	binary.Write(ad, binary.BigEndian, int32(w.Version))
	ad.Write(w.MacAddr)
	binary.Write(ad, binary.BigEndian, int16(w.Flags))
	ad.Write(iv)
	binary.Write(ad, binary.BigEndian, int32(w.DataVersion))
	binary.Write(ad, binary.BigEndian, int32(w.DataLength))
	return ad.Bytes()
}

func decryptGCM(payload, iv []byte, key string, w *InformWrapper) ([]byte, error) {
	block, err := decodeHexKey(key)
	if err != nil {
		return nil, err
	}

	mode, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return nil, err
	}

	_, err = mode.Open(payload[:0], iv, payload, buildAuthData(w, iv))
	if err != nil {
		return nil, err
	}

	// The last block always seems to be garbage, maybe it's padding or
	// something else. I have not looked carefully at it.
	return payload[:len(payload)-aes.BlockSize], nil
}

func decryptCBC(payload, iv []byte, key string) ([]byte, error) {
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
