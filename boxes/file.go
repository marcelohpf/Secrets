package boxes

import (
  "os"
  "log"
  "io/ioutil"
)

func ReadFileBox(path string) string {
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

func WriteFileBox(path, content string) {
  var logger = log.New(os.Stdout, "crypto: ", log.Lshortfile)
  logger.Println("parte content to bytes")
  byte_content := []byte(content)

  logger.Println("Write bytes in file")
  err := ioutil.WriteFile(path, byte_content, 384 )
  if err != nil {
    logger.Fatal(err.Error())
    panic(err.Error())
  }
}
