package boxes

import (
	"golang.org/x/oauth2"
	"secrets/config"
)

type Vault interface {
	ReadBoxItem() (string, error)
	WriteBoxItem(content string) error
}

type Box struct {
	boxPath  string `json:box_path`
	boxName  string `json:box_Name`
	itemName string `json:item_name`
}

func Builder(boxPath, boxName, itemName string, token *oauth2.Token) Vault {
	switch config.BackendStorage {
	case "gdrive":
		return DriveBox{
			Box{
				boxPath:  boxPath,
				boxName:  boxName,
				itemName: itemName,
			},
			token,
		}
	default:
		return LocalBox{
			boxPath:  boxPath,
			boxName:  boxName,
			itemName: itemName,
		}
	}
}
