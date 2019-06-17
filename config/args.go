package config

var (
	KeyName    string
	BoxKeyName string
	KeyPath    string
	ItemName   string
	BoxName    string
	BoxPath    string
	InFile     string
	OutFile    string
	SizeKey    int
	Verbose    bool
	Debug      bool
)

const DefaultKeyPath string = "~/.config/vault"
const DefaultBoxPath string = "~/vault/boxes"

var TokenFile string = "~/.config/vault/token.json"
var CredentialsFile string = "~/.config/vault/credentials.json"

//var BackendStorage string = "gdrive"
var BackendStorage string = "local"

var Server string = "localhost"
var Port string = "8000"
