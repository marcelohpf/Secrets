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

const DefaultKeyPath string = "/home/marcelohpf/vault"
const DefaultBoxPath string = "/home/marcelohpf/vault"

var TokenFile string = "/home/marcelohpf/.config/vault/token.json"
var CredentialsFile string = "/home/marcelohpf/.config/vault/credentials.json"
var BackendStorage string = "gdrive"
