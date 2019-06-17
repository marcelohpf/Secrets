package crypto

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func GetKey(keyPath, keyName string) (string, error) {
	log.Info("Getting key")
	if keyPath == "" && keyName == "" {
		log.Debug("Argument key not present read from console")
		return readFromConsole()
	} else {
		log.Debug("Argument key present.")

		expandedPath, err := dirExpansion(keyPath)
		if err != nil {
			return "", err
		}

		absPath, err := filepath.Abs(expandedPath)
		if err != nil {
			return "", err
		}

		ext := filepath.Ext(keyName)
		if ext != ".key" || ext == "" {
			keyName += ".key"
			log.Debug("Adding key sufix")
		}

		return readFromFile(filepath.Join(absPath, keyName))
	}
}

func SaveKey(keyPath, keyName, key string) error {

	log.WithFields(log.Fields{"keyPath": keyPath, "keyName": keyName}).Info("Saving key.")
	if keyName != "" && keyPath != "" {

		expandedPath, err := dirExpansion(keyPath)
		if err != nil {
			return err
		}

		absPath, err := filepath.Abs(expandedPath)
		log.Debug("Creating dir ", absPath)
		err = os.MkdirAll(absPath, 0755)

		if err != nil {
			log.Debug(err.Error())
			return err
		}
		cleanPath := filepath.Clean(absPath)

		ext := filepath.Ext(keyName)
		if ext != ".key" {
			keyName += ".key"
		}

		log.Debug("Key is present.")
		err = ioutil.WriteFile(filepath.Join(cleanPath, keyName), []byte(key), 384)

		if err != nil {
			log.Debug(err.Error())
			return err
		}
	} else {
		log.Warn("Invalid path or file name")
		return errors.New("error invalid path or file name")
	}
	log.Info("Key saved with success!")
	return nil
}

func readFromFile(path string) (string, error) {

	log.Debug("Open file descriptor")

	file, err := os.Open(path)
	if err != nil {
		log.Debug(err.Error())
		return "", err
	}
	defer file.Close()

	log.Debug("Reading file content")
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Debug(err.Error())
		return "", err
	}

	log.Info("Key readed from file!")
	return string(content), nil
}

func readFromConsole() (string, error) {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Debug(err.Error())
		return "", err

	}
	password := string(bytePassword)

	log.Info("Key readed from console!")
	return strings.TrimSpace(password), nil
}

func dirExpansion(path string) (string, error) {
	// This is a function to expand to absolute path in Unix System like

	if strings.HasPrefix(path, "~/") {
		user, _ := user.Current()
		dir := user.HomeDir
		return filepath.Join(dir, path[2:]), nil
	}
	return filepath.Abs(path)
}
