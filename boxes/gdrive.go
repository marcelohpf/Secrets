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

// Refresh google auth token
func GRefreshAuth() (*oauth2.Token, error) {
  oauthConfig, err := getConfig(config.CredentialsFile)
  if err != nil {
    return nil, err
  }
  return refreshToken(oauthConfig)
}

// Read a file from google drive
func GReadBoxItem(boxPath, boxName, itemName string) (string, error) {
  log.Info("Reading box item.")
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
    log.Warn("Could not initialize a google client.")
    return "", err
  }

  parentId, err := getDirId(srv, boxPath, boxName)
  if err != nil {
    log.Debug("parent id not found")
    return "", err
  }

  itemId, err := getItemGId(srv, parentId, itemName)
  if err != nil {
    log.Debug("Something wrong happen when try retrieve item google id.")
    return "", err
  }
  return fetchRemoteFile(srv, itemId)
}

// Write the content item into a file in the box path in subdir of box name
func GWriteBoxItem(boxPath, boxName, itemName, content string) error {

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

  parentId, err := ensureDirs(srv, boxPath, boxName)

  if err != nil {
    return err
  }

  log.Debug("Creating file under ", parentId)

  return upsert(srv, parentId, itemName, content)
}

func fetchRemoteFile(service *drive.Service, itemId string) (string, error) {
  log.Info("Fetching remote file", itemId)

  http, err := service.Files.Get(itemId).Download()
  if err != nil {
    log.Debug("http drive file retrive", err)
    return "", err
  }
  defer http.Body.Close()

  buff := new(bytes.Buffer)
  buff.ReadFrom(http.Body)
  content := buff.String()

  log.Info("File fetched with success")
  return content, nil
}

func getDirId(service *drive.Service, boxPath, boxName string) (string, error) {
  finalPath := boxPath + "/" + boxName
  finalPath = strings.Trim(finalPath, "/")

  // there is no boxPath or boxName set use root dir
  if finalPath == "" {
    log.Info("Final path is a empty string, using the root dir")
    return ROOTID, nil
  }

  paths := strings.Split(finalPath, "/")

  parentId := ROOTID
  for _, path := range paths {
    log.Debug("Get dir id for ", parentId, " / ", path)
    id, err := getItemGId(service, parentId, path)
    if err != nil {
      return "", err
    }
    parentId = id
  }
  log.Debug("The id ", finalPath, " (", parentId, ")")
  return parentId, nil
}


func ensureDirs(service *drive.Service, boxPath, boxName string) (string, error) {
  finalPath := boxPath + "/" + boxName
  finalPath = strings.Trim(finalPath, "/")

  // there is no boxPath or boxName set use root dir
  if finalPath == "" {
    log.Info("Final path is a empty string, using the root dir")
    return ROOTID, nil
  }

  paths := strings.Split(finalPath, "/")

  parentId := ROOTID
  var subPath int
  for subPath = 0; subPath < len(paths); subPath++ {
    log.Debug("searching id for ", paths[subPath])
    id, err := getItemGId(service, parentId, paths[subPath])
    if err == notFound {
      log.Warn("Subpath does not exists, it will try to create", paths[subPath:])
      break
    } else if err != nil {
      return "", err
    }
    parentId = id
  }

  if subPath == len(paths) {
    log.Debug("Subpath at the end ", parentId)
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

  return parentId, nil
}

func createDir(service *drive.Service, parentId, name string) (*drive.File, error) {
  log.Debug("Create dir for parent:name ", parentId, ":", name)
  d := &drive.File{
    Name:     name,
    MimeType: "application/vnd.google-apps.folder",
    Parents:  []string{parentId},
  }

  file, err := service.Files.Create(d).Do()

  if err != nil {
    log.Println("Could not create dir: " + err.Error())
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
    log.Debug("Update file", file.Id)
  } else { // some other error happens, try to create new the file
    file, err := createFile(service, itemName, content, parentId)
    if err != nil {
      return err
    }
    log.Debug("Created file: ", file.Parents, " / ", file.Id)
  }
  return nil
}

// Retrieve a google id for a given item name
func getItemGId(service *drive.Service, parentId, itemName string) (string, error) {
  log.Debug("Finding file on drive ", parentId, " ", itemName)
  flCall := service.Files.List().PageSize(10).
    Fields("nextPageToken, files(id, name)").
    Q("name = '"+ itemName + "' and '" + parentId + "' in parents")

  r, err := flCall.Do()

  if  err != nil {
    log.Debug("Unable to get id request")
    return "", err
  }

  log.Debug("Found files ", len(r.Files))

  switch len(r.Files) {
  case 0:
    log.Warn("Item not found on box. ", parentId, " ", itemName)
    return "", notFound
  case 1:
    log.Debug("Found the file ", parentId, " / ", r.Files[0].Id)
    return r.Files[0].Id, nil
  default: // should use the last one ordered by modification date
    for _, f := range r.Files {
      log.Debug("Files: ", parentId, " / ", f.Id, " (", f.Name, ")")
    }
    return "", foundMany
  }
}

// Retrieve a token, saves the token or load from cache.
func getToken(oauthConfig *oauth2.Config) (*oauth2.Token, error) {
  token, err := tokenFromFile(config.TokenFile)
  if err != nil {
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
  WriteIntoFile(config.TokenFile, string(tokenJson))
  return token, nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(oauthConfig *oauth2.Config) *oauth2.Token {
  authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
  fmt.Printf("Auhtorization URL\n%v\nAuthorization code:\n", authURL)

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
    log.Warn("Token not retrieved from file")
    return nil, err
  }
  tok := &oauth2.Token{}
  err = json.Unmarshal([]byte(jsonToken), tok)
  if err != nil {
    log.Warn("Could not Unmarshal json token")
    return nil, err
  }

  if tok.Valid() {
    log.Debug("Token from file is ok!")
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
    log.Warn("Unable to parse client secret file to config")
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
    log.Debug("Could not update file: " + err.Error())
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
    log.Debug("Could not create file: " + err.Error())
    return nil, err
  }

  return file, nil
}
