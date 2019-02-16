package config

import (
  "flag"
)

var (
  Key string
  Text string
  TextPath string
  CiphertextPath string
  Encrypt bool
  Decrypt bool
)

func init(){
  flag.StringVar(&Key, "k", "", "Secret key.")
  flag.StringVar(&Text, "t", "", "Texto to cipher or decipher")
  flag.StringVar(&TextPath, "tp", "", "Alternative to use a path to file.")
  flag.StringVar(&CiphertextPath, "cp", "", "Select where the cipher content will be writen.")
  flag.BoolVar(&Encrypt, "e", false, "Operation of encrypt a file.")
  flag.BoolVar(&Decrypt, "d", false, "Operation of decrypt a file.")
  flag.Parse()
}
