package dri

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"google.golang.org/api/googleapi"

	"golang.org/x/oauth2"

	"parsecsreach.com/backuptool/conf"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(cfg *conf.Config) (*http.Client, error) {

	if cfg.Drive.Token == nil {
		tok, err := getTokenFromWeb(cfg.Drive.OauthConfig)
		if err != nil {
			return nil, err
		}
		cfg.Drive.Token = tok
		conf.WriteConfig(cfg)
	}

	return cfg.Drive.OauthConfig.Client(context.Background(), cfg.Drive.Token), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok, nil
}

/*
GetClient returns a new google drive client if a token is not saved this will ask for one.
*/
func GetClient(config *conf.Config) (*drive.Service, error) {
	client, err := getClient(config)
	if err != nil {
		return nil, err
	}
	srv, err := drive.New(client)
	if err != nil {
		return nil, err
	}
	return srv, err
}

/*
FindOrCreateFolder will find a folder in a google drive and return the id. if it doesn't exist it will create it first.
*/
func FindOrCreateFolder(conn *drive.Service, name string) (string, error) {

	id, err := FindID(conn, name, "application/vnd.google-apps.folder")
	if err != nil {
		return "", err
	}

	if id == "" {
		id, err = createFolder(conn, name)
		if err != nil {
			return "", err
		}
	}

	return id, nil
}

/*
FindID returns the file id for a given name and mimeType or blank if it cant find it
*/
func FindID(conn *drive.Service, name string, mimeType string) (string, error) {

	filelist, err := conn.Files.List().
		Q(fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder'", name)).
		IncludeTeamDriveItems(false).Do()

	if err != nil {
		return "", err
	}

	for _, f := range filelist.Files {
		return f.Id, nil
	}

	return "", nil
}

func createFolder(conn *drive.Service, name string) (string, error) {
	f := drive.File{
		MimeType: "application/vnd.google-apps.folder",
		Name:     name,
	}
	r, err := conn.Files.Create(&f).Fields(googleapi.Field("id")).Do()
	if err != nil {
		return "", err
	}
	return r.Id, nil
}

/*
Upload loads the content into google drive
*/
func Upload(conn *drive.Service, name string, mimeType string, parent string, size int64, content io.ReaderAt) (string, error) {
	f := drive.File{
		Name:     name,
		MimeType: mimeType,
	}

	if parent != "" {
		f.Parents = []string{parent}
	}

	r, err := conn.Files.Create(&f).ResumableMedia(context.Background(), content, size, mimeType).Do()
	if err != nil {
		return "", err
	}
	return r.Id, nil
}
