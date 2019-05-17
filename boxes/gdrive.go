package boxes

import (
  "secrets/config"
  "encoding/json"
  "fmt"
  log "github.com/sirupsen/logrus"
  "errors"
  "strings"
  "bytes"

  "golang.org/x/net/context"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/drive/v3"
  "google.golang.org/api/option"
)

var notFound = errors.New("Record not found.")
var foundMany = errors.New("Found more then one record.")
const ROOTID string = "root"

type DriveBox Box

// Refresh google auth token
func GRefreshAuth() (*oauth2.Token, error) {
  oauthConfig, err := getConfig(config.CredentialsFile)
  if err != nil {
    return nil, err
  }
  return refreshToken(oauthConfig)
}

// Read a file from google drive
func (box DriveBox) ReadBoxItem() (string, error) {
  log.WithFields(log.Fields{
    "boxPath": box.boxPath,
    "boxName": box.boxName,
    "itemName": box.itemName,
  }).Info("Reading google drive secret box ")
  gconfig, err := getConfig(config.CredentialsFile)
  if err != nil {
    return "", err
  }

  token, err := getToken(gconfig)
  if err != nil {
    return "", err
  }

  ctx := context.Background()
  srv, err := drive.NewService(ctx, option.WithTokenSource(gconfig.TokenSource(ctx, token)))
  if err != nil {
    log.WithFields(log.Fields{
      "boxPath": box.boxPath,
      "boxName": box.boxName,
      "itemName": box.itemName,
    }).Warn("Could not initialize a google client.")
    return "", err
  }

  parentId, err := getDirId(srv, box)
  if err != nil {
    log.WithFields(log.Fields{
      "boxPath": box.boxPath,
      "boxName": box.boxName,
      "itemName": box.itemName,
    }).Warn("Dir id not found")
    return "", err
  }

  if box.itemName == "" {
    return "", errors.New("Item name not defined")
  }

  item := addSufix(box.itemName)
  itemId, err := getItemGId(srv, parentId, item)
  if err != nil {
    log.WithFields(log.Fields{
      "parentId": parentId,
      "itemName": item,
    }).Debug("Error happen when try get the google id ")
    return "", err
  }
  return fetchRemoteFile(srv, itemId)
}

// Write the content item into a file in the box path in subdir of box name
func (box DriveBox) WriteBoxItem(content string) error {

  gconfig, err := getConfig(config.CredentialsFile)
  if err != nil {
    panic("error to write")
  }

  token, err := getToken(gconfig)
  ctx := context.Background()
  srv, err := drive.NewService(ctx, option.WithTokenSource(gconfig.TokenSource(ctx, token)))

  if err != nil {
    log.Fatal(err.Error())
  }

  parentId, err := ensureDirs(srv, box)

  if err != nil {
    return err
  }

  log.WithFields(log.Fields{
    "parentId": parentId,
  }).Debug("Creating file")

  if box.itemName == "" {
    return errors.New("Item name not defined")
  }

  item := addSufix(box.itemName)
  return upsert(srv, parentId, item, content)
}

func fetchRemoteFile(service *drive.Service, itemId string) (string, error) {

  http, err := service.Files.Get(itemId).Download()
  if err != nil {
    log.WithFields(log.Fields{
      "itemName": itemId,
    }).Debug("HTTP drive file retrive error ", itemId)
    return "", err
  }
  defer http.Body.Close()

  buff := new(bytes.Buffer)
  buff.ReadFrom(http.Body)
  content := buff.String()

  log.WithFields(log.Fields{
    "itemName": itemId,
  }).Info("File fetched from google drive")
  return content, nil
}

func getDirId(service *drive.Service, box DriveBox) (string, error) {

  path, _ := gdirExpansion(box.boxPath)
  finalPath :=  path + "/" + box.boxName
  finalPath = strings.Trim(finalPath, "/")

  // there is no boxPath or boxName set use root dir
  if finalPath == "" {
    log.WithFields(log.Fields{
      "boxPath": box.boxPath,
      "boxName": box.boxName,
    }).Info("Final path is a empty string, using the root dir")
    return ROOTID, nil
  }

  paths := strings.Split(finalPath, "/")

  parentId := ROOTID
  for _, path := range paths {
    log.WithFields(log.Fields{
      "parentId": parentId,
      "itemName": path,
    }).Debug("Get dir id")
    id, err := getItemGId(service, parentId, path)
    if err != nil {
      return "", err
    }
    parentId = id
  }
  log.WithFields(log.Fields{
    "parentId": parentId,
    "itemName": finalPath,
  }).Debug("Found the final id")
  return parentId, nil
}


func ensureDirs(service *drive.Service, box DriveBox) (string, error) {
  path, _ := gdirExpansion(box.boxPath)
  finalPath := path + "/" + box.boxName
  finalPath = strings.Trim(finalPath, "/")

  // there is no boxPath or boxName set use root dir
  if finalPath == "" {
    log.WithFields(log.Fields{
      "boxPath": box.boxPath,
      "boxName": box.boxName,
    }).Info("Final path is a empty string, using the root dir")
    return ROOTID, nil
  }

  paths := strings.Split(finalPath, "/")

  parentId := ROOTID
  var subPath int
  for subPath = 0; subPath < len(paths); subPath++ {
    log.WithFields(log.Fields{
      "boxPath": box.boxPath,
      "boxName": box.boxName,
    }).Debug("Searching id ", paths[subPath])
    id, err := getItemGId(service, parentId, paths[subPath])
    if err == notFound {
      log.WithFields(log.Fields{
        "boxPath": box.boxPath,
        "boxName": box.boxName,
      }).Warn("Subpath does not exists, it will try to create ", strings.Join(paths[subPath:], "/"))
      break
    } else if err != nil {
      return "", err
    }
    parentId = id
  }

  if subPath == len(paths) {
    log.WithFields(log.Fields{
      "parentId": parentId}).Debug("Subpath final parent id ", parentId)
    return parentId, nil
  } else {
    // it creates missing dirs or just return the parent id
    return createDirs(service, parentId, paths[subPath:])
  }
}

func createDirs(service *drive.Service, parentId string, names []string) (string, error) {

  for _, name := range names {
    dir, err := createDir(service, parentId, name)
    if err != nil {
      return "", err
    }
    parentId = dir.Id
  }

  log.WithFields(log.Fields{
    "parentId": parentId,
    "itemName": strings.Join(names, "/"),
  }).Info("Created missing paths")

  return parentId, nil
}

func createDir(service *drive.Service, parentId, name string) (*drive.File, error) {
  log.WithFields(log.Fields{
    "parentId": parentId,
    "itemName": name,
  }).Debug("Creating folder ")
  d := &drive.File{
    Name:     name,
    MimeType: "application/vnd.google-apps.folder",
    Parents:  []string{parentId},
  }

  file, err := service.Files.Create(d).Do()

  if err != nil {
    log.WithFields(log.Fields{
      "parentId": parentId,
      "itemName": name,
    }).Debug("Could not create dir")
    return nil, err
  }

  return file, nil
}

// insert or update a item
func upsert(service *drive.Service, parentId, itemName, content string) error {
  gid, err := getItemGId(service, parentId, itemName)

  // found item do a update
  if err == nil {
    file, err := updateFile(service, gid, itemName, content)
    if err != nil {
      return err
    }
    log.WithFields(log.Fields{
      "parentId": parentId,
      "itemName": file.Id,
      "name": file.Name,
    }).Info("Update file")
  } else { // some other error happens, try to create new the file
    file, err := createFile(service, itemName, content, parentId)
    if err != nil {
      return err
    }
    log.WithFields(log.Fields{
      "parentId": parentId,
      "itemName": file.Id,
      "name": file.Name,
    }).Info("Created file ")
  }
  return nil
}

// Retrieve a google id for a given item name
func getItemGId(service *drive.Service, parentId, itemName string) (string, error) {
  log.WithFields(log.Fields{
    "parentId": parentId,
    "itemName": itemName,
  }).Debug("Finding file on drive")
  flCall := service.Files.List().PageSize(10).
    Fields("files(id, name)").
    Q("name = '"+ itemName + "' and '" + parentId + "' in parents")

  r, err := flCall.Do()

  if  err != nil {
    log.WithFields(log.Fields{
      "parentId": parentId,
      "itemName": itemName,
    }).Debug("Unable to get id request")
    return "", err
  }

  log.WithFields(log.Fields{
    "parentId": parentId,
    "itemName": itemName,
  }).Debug("Found files ", len(r.Files))

  switch len(r.Files) {
  case 0:
    log.WithFields(log.Fields{
      "parentId": parentId,
      "itemName": itemName,
    }).Warn("Item not found on box. ")
    return "", notFound
  case 1:
    log.WithFields(log.Fields{
      "parentId": parentId,
      "itemName": r.Files[0].Id,
    }).Debug("Found the file ")
    return r.Files[0].Id, nil
  default: // should use the last one ordered by modification date
    for _, f := range r.Files {
      log.WithFields(log.Fields{
        "parentId": parentId,
        "itemName": f.Name,
      }).Debug("Files: ", parentId, " / ", f.Id, " (", f.Name, ")")
    }
    return "", foundMany
  }
}

// Retrieve a token, saves the token or load from cache.
func getToken(oauthConfig *oauth2.Config) (*oauth2.Token, error) {
  token, err := tokenFromFile(config.TokenFile)
  if err != nil {
    log.Error(err.Error())
    return refreshToken(oauthConfig)
  }
  return token, nil
}

func refreshToken(oauthConfig *oauth2.Config) (*oauth2.Token, error) {
  token := getTokenFromWeb(oauthConfig)
  log.Debug("Saving credential file.")
  tokenJson, err := json.Marshal(token)
  if err != nil {
    log.Debug("Could not parse the token")
    return nil, err
  }
  err = WriteIntoFile(config.TokenFile, string(tokenJson))
  return token, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(oauthConfig *oauth2.Config) *oauth2.Token {
  authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
  fmt.Printf("Authorization URL\n%v\nAuthorization code:\n", authURL)

  var authCode string
  if _, err := fmt.Scan(&authCode); err != nil {
    log.Fatal("Unable to read authorization code %v", err)
  }

  tok, err := oauthConfig.Exchange(context.TODO(), authCode)
  if err != nil {
    log.Fatal("Unable to retrieve token from web %v", err)
  }
  return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
  jsonToken, err := ReadFromFile(file)
  if err != nil {
    log.WithFields(log.Fields{
      "file": file,
    }).Warn("Token not retrieved from file")
    return nil, err
  }
  tok := &oauth2.Token{}
  err = json.Unmarshal([]byte(jsonToken), tok)
  if err != nil {
    log.WithFields(log.Fields{
      "file": file,
    }).Warn("Could not Unmarshal json token")
    return nil, err
  }

  if tok.Valid() {
    log.WithFields(log.Fields{
      "file": file,
    }).Debug("Token from file is ok!")
    return tok, nil
  } else {
    return nil, errors.New("Token is invalid or expired")
  }

}

func getConfig(credentialsFile string) (*oauth2.Config, error) {
  b, err := ReadFromFile(credentialsFile)
  if err != nil {
    return nil, err
  }

  // If modifying these scopes, delete your previously saved token.json.
  config, err := google.ConfigFromJSON([]byte(b), drive.DriveFileScope)
  if err != nil {
    log.WithFields(log.Fields{
      "file": credentialsFile,
    }).Warn("Unable to parse client secret file to config")
    return nil, err
  }
  return config, nil
}

// update files not change the parent ID
func updateFile(service *drive.Service, gid, name, content string) (*drive.File, error) {
  f := &drive.File{
    Name:     name,
  }

  ioContent := strings.NewReader(content)
  file, err := service.Files.Update(gid, f).Media(ioContent).Do()

  if err != nil {
    log.WithFields(log.Fields{
      "itemName": gid,
      "name": name,
    }).Debug("Could not update file: " + err.Error())
    return nil, err
  }

  return file, nil
}

func createFile(service *drive.Service, name, content, parentId string) (*drive.File, error) {
  f := &drive.File{
    Name:     name,
    Parents: []string{parentId},
  }

  ioContent := strings.NewReader(content)
  file, err := service.Files.Create(f).Media(ioContent).Do()

  if err != nil {
    log.WithFields(log.Fields{
      "parentId": parentId,
      "itemName": name,
    }).Debug("Could not create file")
    return nil, err
  }

  return file, nil
}

func gdirExpansion(path string) (string, error) {
  return strings.Trim(path, "~/"), nil
}
