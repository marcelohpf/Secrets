package boxes

import (
  "os"
  log "github.com/sirupsen/logrus"
  "io/ioutil"
  "path/filepath"
  "errors"
)

func ReadBoxItem(boxPath, boxName, itemName string) (string, error) {
  log.Info("Reading item from box")
  box, err := mountBox(boxPath, boxName)
  if err != nil {
    return "", err
  }
  item, err := mountItem(box, itemName)
  if err != nil {
    return "", err
  }
  return ReadFromFile(item)
}

func WriteBoxItem(boxPath, boxName, itemName, content string) error {
  log.Info("Writing item into box")

  box, err := mountBox(boxPath, boxName)

  if err != nil {
    return err
  }

  item, err := mountItem(box, itemName)

  if err != nil {
    return err
  }
  return WriteIntoFile(item, content)
}

func ReadFromFile(path string) (string, error) {
  absPath, err := filepath.Abs(path)
  if err != nil {
    return "", err
  }

  log.Debug("Open file descriptor.", absPath)
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

func WriteIntoFile(path, content string) error {
  byteContent := []byte(content)

  absPath, err := filepath.Abs(path)
  if err != nil {
    return err
  }

  dir := filepath.Dir(absPath)
  err = os.MkdirAll(dir, 0777)

  if err != nil {
    log.Debug(err.Error())
    return err
  }

  log.Debug("Write bytes in file ")
  err = ioutil.WriteFile(path, byteContent, 384)
  if err != nil {
    log.Debug(err)
    return err
  }
  log.Debug("Bytes writed with success")
  return nil
}

func mountBox(boxPath, boxName string) (string, error) {
  cleanPath := filepath.Clean(boxPath)

  absPath, err := filepath.Abs(cleanPath)

  if err != nil {
    return "", err
  }

  if boxName != "" {
    absPath = filepath.Join(absPath, boxName)
  } else {
    log.Debug("No box was selected")
  }

  return absPath, nil
}

func mountItem(box, itemName string) (string, error) {

  if itemName == "" {
    log.Debug("Empty item name")
    return "", errors.New("No item was selected.")
  }
  log.Debug("Non item name", itemName)

  item := addSufix(filepath.Join(box, itemName))

  return item, nil
}

func addSufix(box string) (string) {
  ext := filepath.Ext(box)
  if ext == "" || ext != ".vlt" {
    return box + ".vlt"
  }
  return box
}
