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

  itemGId, err := getItemGId(boxPath, boxName, itemName, srv)
  if err != nil {
    log.Warn("GID not retrieved")
    return "", err
  }

  http, err := srv.Files.Get(itemGId).Download()
  if err != nil {
    log.Debug("http drive file retrive", err)
    return "", err
  }

  defer http.Body.Close()

  buff := new(bytes.Buffer)
  buff.ReadFrom(http.Body)
  content := buff.String()

  log.Debug(boxPath, boxName, itemName)
  return content, nil
}

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
  file, err := createFile(srv, itemName, content, "root")
  if err != nil {
    return err
  }
  log.Info("Create a file with id ", file.Id)
  return nil
}

func getItemGId(boxPath, boxName, itemName string, service *drive.Service) (string, error) {
  log.Debug("Finding file on drive.")
  r, err := service.Files.List().PageSize(1).
    Fields("nextPageToken, files(id, name)").
    Q("name = '"+ itemName + "'").Do()

  if err != nil {
    log.Debug("Unable to get id request")
    return "", err
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
func getToken(oauthConfig *oauth2.Config) (*oauth2.Token, error) {
  // The file token.json stores the user's access and refresh tokens, and is
  // created automatically when the authorization flow completes for the first
  // time.
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

func createFile(service *drive.Service, name, content, parentId string) (*drive.File, error) {

   f := &drive.File{
      Name:     name,
   }

   ioContent := strings.NewReader(content)
   file, err := service.Files.Create(f).Media(ioContent).Do()

   if err != nil {
      log.Debug("Could not create file: " + err.Error())
      return nil, err
   }

   return file, nil
 }
