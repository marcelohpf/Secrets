package config

var (
  KeyPath string
  KeyName string
  ItemName string
  BoxName string
  BoxPath string
  InFile string
  OutFile string
  SizeKey int
  Verbose bool
  Debug bool
)

const DefaultKeyPath string = "vault"
const DefaultBoxPath string = "vault"

var TokenFile string = "~/.config/vault/token.json"
var CredentialsFile string = "~/.config/vault/credentials.json"
var BackendStorage string = "gdrive"
