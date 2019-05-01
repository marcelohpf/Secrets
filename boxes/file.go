package boxes

import (
  "os"
  log "github.com/sirupsen/logrus"
  "io/ioutil"
)


func ReadBoxItem(boxPath, boxName, itemName string) string {
  log.Info("Reading item from box")
  if boxPath != "" && itemName != "" {
    return ReadFromFile(boxPath + "/" + boxName + "/" + itemName)
  }
  return ""
}

func WriteBoxItem(boxPath, boxName, itemName, content string) {
  log.Info("Writing item into box")
  if boxPath != "" && itemName != "" {
    WriteIntoFile(boxPath + "/" + boxName + "/" + itemName, content)
  }
}

func ReadFromFile(path string) string {
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

  log.Debug("File content readed from file!")
  return string(content)
}

func WriteIntoFile(path, content string) {
  log.Debug("parse content to bytes")
  byteContent := []byte(content)

  log.Debug("Write bytes in file")
  err := ioutil.WriteFile(path, byteContent, 384)
  if err != nil {
    log.Fatal(err.Error())
    panic(err.Error())
  }
  log.Debug("Bytes writed with success")
}
