package crypto

import (
  "fmt"
  "os"
  "log"
  "io/ioutil"
  "syscall"
  "strings"

  "vault/config"

  "golang.org/x/crypto/ssh/terminal"
)

func GetKey() string {
  var logger = log.New(os.Stdout, "crypto ", log.Lshortfile)
  logger.Println("Getting key")
  if config.Key != "" {
    logger.Println("Argument key present.")
    return read_from_file(config.Key)
  } else {
    logger.Println("Argument key not present")
    return read_from_console()
  }
}

func SaveKey(key string) {
  var logger = log.New(os.Stdout, "crypto ", log.Lshortfile)
  logger.Println("Saving key")
  if config.Key != "" {
    logger.Println("Key path present.")
    err := ioutil.WriteFile(config.Key, []byte(key), 384)
    if err != nil {
      logger.Fatal(err.Error())
      panic(err.Error())
    }
  }
  fmt.Println(key)
}

func read_from_file(path string) string {
  var logger = log.New(os.Stdout, "crypto ", log.Lshortfile)

  logger.Println("Open file descriptor")
  file, err := os.Open(path)
  if err != nil {
    logger.Fatal(err.Error())
    panic(err.Error())
  }
  defer file.Close()

  logger.Println("Read file content")
  content, err := ioutil.ReadAll(file)
  if err != nil {
    logger.Fatal(err.Error())
    panic(err.Error())
  }

  return string(content)
}

func read_from_console() string {

  var logger = log.New(os.Stdout, "crypto ", log.Lshortfile)
  fmt.Print("Enter Password: ")
  bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
  if err != nil {
    logger.Fatal(err.Error())
    panic(err.Error())

  }
  password := string(bytePassword)

  return strings.TrimSpace(password)
}
