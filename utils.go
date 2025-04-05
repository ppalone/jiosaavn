package jiosaavn

import (
	"crypto/des"
	"encoding/base64"
	"fmt"
)

func generateMediaURL(encryptedMediaURL string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encryptedMediaURL)
	if err != nil {
		return "", err
	}

	block, err := des.NewCipher([]byte("38346591"))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if len(decodedBytes)%blockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	decrypted := make([]byte, len(decodedBytes))
	for start := 0; start < len(decodedBytes); start += blockSize {
		block.Decrypt(decrypted[start:start+blockSize], decodedBytes[start:start+blockSize])
	}

	return string(stripPKCS5Padding(decrypted)), nil
}

func stripPKCS5Padding(data []byte) []byte {
	paddingLen := int(data[len(data)-1])
	return data[:len(data)-paddingLen]
}
