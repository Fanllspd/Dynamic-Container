package utils

import (
	"encoding/base64"
	"os"
)

func LoadBase64EncodedFile(filePath string) ([]byte, error) {
	encoded, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	decoded, err := base64.StdEncoding.DecodeString(string(encoded))
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
