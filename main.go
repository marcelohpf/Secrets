package main

import (
  log "github.com/sirupsen/logrus"
  "secrets/cmd"
)

func main(){
  log.Info("Personal Vault Service")
  cmd.Execute()
}
