package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"secrets/config"
)

func GenerateKey() (string, error) {
	log.Info("Generating a new key.")
	log.Debug("Validate Key Size")
	if config.SizeKey != 16 && config.SizeKey != 32 {
		log.Fatal("Key should have size of 16 or 32 bytes.")
		return "", errors.New("The key is from a unexpected size")
	}

	key, err := generate(config.SizeKey)
	if err != nil {
		return "", err
	}
	log.Debug("Encoding generated key")
	encodedKey := base64.StdEncoding.EncodeToString(key)
	log.Info("New key generated with success!")
	return encodedKey, nil
}

func generate(size int) ([]byte, error) {
	log.Info("Generating stream")
	data := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		log.Debug(err.Error())
		return []byte{}, err
	}
	return data, nil
}
