package boxes

import (
  "os"
  log "github.com/sirupsen/logrus"
  "io/ioutil"
)


func ReadBoxItem(boxPath, boxName, itemName string) (string, error) {
  log.Info("Reading item from box")
  if boxPath != "" && itemName != "" {
    return ReadFromFile(boxPath + "/" + boxName + "/" + itemName)
  }
  return "", nil
}

func WriteBoxItem(boxPath, boxName, itemName, content string) {
  log.Info("Writing item into box")
  if boxPath != "" && itemName != "" {
    WriteIntoFile(boxPath + "/" + boxName + "/" + itemName, content)
  }
}

func ReadFromFile(path string) (string, error) {
  log.Debug("Open file descriptor.", path)
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    log.Debug("Problem to read file", path )
    return "", err
  }

  log.Debug("Reading file content.")
  content, err := ioutil.ReadAll(file)
  if err != nil {
    log.Debug("Fail to read content of file.", path)
    return "", err
  }

  log.Debug("File content readed from file!")
  return string(content), nil
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
