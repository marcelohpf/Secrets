package boxes

import (
  "secrets/config"
)

type Vault interface {
  ReadBoxItem() (string, error)
  WriteBoxItem(content string) error
}

type Box struct {
  boxPath string `json:box_path`
  boxName string `json:box_Name`
  itemName string `json:item_name`
}

func Builder(boxPath, boxName, itemName string) Vault {
  switch config.BackendStorage {
    case "gdrive": 
      return DriveBox{
        boxPath: config.BoxPath,
        boxName: config.BoxName,
        itemName: config.ItemName,
      }
    default:
      return LocalBox{
        boxPath: config.BoxPath,
        boxName: config.BoxName,
        itemName: config.ItemName,
      }
  }
}

