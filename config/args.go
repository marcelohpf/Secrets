package config

var (
	KeyPath  string
	KeyName  string
	ItemName string
	BoxName  string
	BoxPath  string
	InFile   string
	OutFile  string
	SizeKey  int
	Verbose  bool
	Debug    bool
)

const DefaultKeyPath string = "~/vault/boxes"
const DefaultBoxPath string = "~/vault/boxes"

var TokenFile string = "~/.config/vault/token.json"
var CredentialsFile string = "~/.config/vault/credentials.json"
var BackendStorage string = "gdrive"

var Server string = "localhost"
var Port string = "8000"
