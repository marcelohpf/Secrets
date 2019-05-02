package boxes

import (
  "secrets/config"
  "encoding/json"
  "fmt"
  log "github.com/sirupsen/logrus"
  "net/http"
  "errors"

  "golang.org/x/net/context"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/drive/v3"
)

func GReadBoxItem(boxPath, boxName, itemName string) string {
  log.Info("Reading box item.")
  gconfig := getConfig(config.CredentialsFile)
  client := getClient(gconfig)

  srv, err := drive.New(client)
  if err != nil {
    log.Fatal(err.Error())
    panic("Could not initialize a google client.")
  }

  itemGId, err := getItemGId(boxPath, boxName, itemName, srv)
  if err != nil {
    log.Fatal(err.Error())
    panic("Google id not found.")
  }

  http, err := srv.Files.Get(itemGId).Download()
  if err != nil {
    log.Fatal(err.Error())
    panic("Impossible to download file.")
  }
  log.Info(http)

  log.Debug(boxPath, boxName, itemName)
  return ""
}

func GWriteBoxItem(boxPath, boxName, itemName, content string) {
  gconfig := getConfig(config.CredentialsFile)
  client := getClient(gconfig)

  srv, err := drive.New(client)
  if err != nil {
    log.Fatal(err.Error())
  }

  r, err := srv.Files.List().PageSize(10).
  Fields("nextPageToken, files(id, name)").Do()
  if err != nil {
    log.Fatal(err.Error())
  }
  fmt.Println("Files:")
  if len(r.Files) == 0 {
    fmt.Println("No files found.")
  } else {
    for _, i := range r.Files {
      fmt.Printf("%s (%s)\n", i.Name, i.Id)
    }
  }
}

func getItemGId(boxPath, boxName, itemName string, service *drive.Service) (string, error) {
  log.Debug("Finding file on drive.")
  r, err := service.Files.List().PageSize(1).
    Fields("nextPageToken, files(id, name)").
    Q("name = '"+ itemName + "'").Do()

  if err != nil {
    log.Fatal(err.Error())
  }

  switch len(r.Files) {
  case 0:
    log.Warn("Item not found on box.", boxName)
    return "", errors.New("Record not found")
  case 1:
    return r.Files[0].Id, nil
  default:
    return "", errors.New("Found more then one record.")
  }
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(oauthConfig *oauth2.Config) *http.Client {
  // The file token.json stores the user's access and refresh tokens, and is
  // created automatically when the authorization flow completes for the first
  // time.
  token, err := tokenFromFile(config.TokenFile)
  if err != nil {
    token = getTokenFromWeb(oauthConfig)
    log.Debug("Saving credential file.")
    tokenJson, err := json.Marshal(token)
    if err != nil {
      log.Fatal(err.Error())
      panic("that is bad")
    }
    WriteIntoFile(config.TokenFile, string(tokenJson))
  }
  return oauthConfig.Client(context.Background(), token)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(oauthConfig *oauth2.Config) *oauth2.Token {
  authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
  fmt.Printf("Go to the following link in your browser then type the "+
  "authorization code: \n%v\n", authURL)

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
  tok := &oauth2.Token{}
  jsonToken := ReadFromFile(file)
  err := json.Unmarshal([]byte(jsonToken), tok)
  return tok, err
}

func getConfig(credentialsFile string) *oauth2.Config {
  b := ReadFromFile(credentialsFile)

  // If modifying these scopes, delete your previously saved token.json.
  config, err := google.ConfigFromJSON([]byte(b), drive.DriveMetadataReadonlyScope)
  if err != nil {
    log.Fatal("Unable to parse client secret file to config: %v", err)
  }
  return config
}
