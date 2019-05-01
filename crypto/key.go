package crypto

import (
  "fmt"
  "os"
  log "github.com/sirupsen/logrus"
  "io/ioutil"
  "syscall"
  "strings"

  "secrets/config"

  "golang.org/x/crypto/ssh/terminal"
)

func GetKey(keyPath, keyName string) string {
  log.Info("Getting key")
  if keyPath == "" && keyName == "" {
    log.Debug("Argument key not present read from console")
    return readFromConsole()
  } else {
    log.Debug("Argument key present.")
    return readFromFile(keyPath + "/" + keyName)
  }
}

func SaveKey(keyPath, keyName, key string) {
  log.Info("Saving key.")
  if keyName != "" && keyPath != "" {
    log.Debug("Key is present.")
    err := ioutil.WriteFile(config.KeyPath + "/" + config.KeyName, []byte(key), 384)
    if err != nil {
      log.Fatal(err.Error())
      panic(err.Error())
    }
  } else {
    log.Error("Invalid path or file name")
    panic("error invalid path or file name")
  }
  log.Info("Key saved with success!")
}

func readFromFile(path string) string {
  log.Debug("Open file descriptor")
  file, err := os.Open(path)
  if err != nil {
    log.Fatal(err.Error())
    panic(err.Error())
  }
  defer file.Close()

  log.Debug("Reading file content")
  content, err := ioutil.ReadAll(file)
  if err != nil {
    log.Fatal(err.Error())
    panic(err.Error())
  }

  log.Info("Key readed from file!")
  return string(content)
}

func readFromConsole() string {
  fmt.Print("Enter Password: ")
  bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
  if err != nil {
    log.Fatal(err.Error())
    panic(err.Error())

  }
  password := string(bytePassword)

  log.Info("Key readed from console!")
  return strings.TrimSpace(password)
}
