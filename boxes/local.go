package boxes

import (
  "os"
  "os/user"
  "strings"
  log "github.com/sirupsen/logrus"
  "io/ioutil"
  "path/filepath"
  "errors"
)

func ReadBoxItem(boxPath, boxName, itemName string) (string, error) {
  log.WithFields(log.Fields{
    "boxPath": boxPath,
    "boxName": boxName,
    "itemName": itemName,
  }).Info("Reading item from box")
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
  log.WithFields(log.Fields{
    "boxPath": boxPath,
    "boxName": boxName,
    "itemName": itemName,
  }).Info("Writing item into box")

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

  absPath, err := dirExpansion(path)

  if err != nil {
    return "", err
  }

  log.WithFields(log.Fields{
    "file": absPath,
  }).Debug("Open file descriptor.",)
  file, err := os.Open(absPath)
  defer file.Close()
  if err != nil {
    log.WithFields(log.Fields{
      "file": absPath,
    }).Debug("Problem to read file",)
    return "", err
  }

  log.WithFields(log.Fields{
    "file": absPath,
  }).Debug("Reading file content.")
  content, err := ioutil.ReadAll(file)
  if err != nil {
    log.WithFields(log.Fields{
      "file": absPath,
    }).Debug("Fail to read content of file.")
    return "", err
  }

  log.WithFields(log.Fields{
    "file": absPath,
  }).Debug("File content readed from file!")
  return string(content), nil
}

func WriteIntoFile(path, content string) error {
  byteContent := []byte(content)

  absPath, err := dirExpansion(path)
  if err != nil {
    return err
  }

  dir := filepath.Dir(absPath)
  err = os.MkdirAll(dir, 0777)

  if err != nil {
    log.Debug(err.Error())
    return err
  }

  log.WithFields(log.Fields{
    "file": absPath,
  }).Debug("Write bytes in file ")
  err = ioutil.WriteFile(path, byteContent, 384)
  if err != nil {
    log.WithFields(log.Fields{
      "file": absPath,
    }).Debug(err.Error())
    return err
  }
  log.WithFields(log.Fields{
    "file": absPath,
  }).Debug("Bytes writed with success")
  return nil
}

func mountBox(boxPath, boxName string) (string, error) {

  absPath, err := dirExpansion(boxPath)

  if err != nil {
    return "", err
  }

  if boxName != "" {
    absPath = filepath.Join(absPath, boxName)
  } else {
    log.WithFields(log.Fields{
      "boxPath": boxPath,
      "boxName": boxName,
    }).Debug("No box was selected")
  }

  return absPath, nil
}

func mountItem(box, itemName string) (string, error) {

  if itemName == "" {
    log.Debug("Empty item name")
    return "", errors.New("No item was selected.")
  }
  log.WithFields(log.Fields{
    "parentId": box,
    "itemName": itemName,
  }).Debug("Non item name", itemName)

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

func dirExpansion(path string) (string, error) {
  // This is a function to expand to absolute path in Unix System like

  if strings.HasPrefix(path, "~/") {
    user, _ := user.Current()
    dir := user.HomeDir
    return filepath.Join(dir, path[2:]), nil
  }
  return filepath.Abs(path)

}
