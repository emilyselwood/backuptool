package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"parsecsreach.com/backuptool/conf"
	"parsecsreach.com/backuptool/dri"
)

/*
Init implements a command that sets up the basic required config
*/
type Init struct {
}

/*
Run Sets up the initial configuration. Will error if the configuration can already be loaded
*/
func (i Init) Run() {

	var force bool

	if len(os.Args) > 2 {
		mySet := flag.NewFlagSet("", flag.ExitOnError)
		mySet.BoolVar(&force, "f", false, "force recreation even if file already exists")
		mySet.Parse(os.Args[2:])
	} else {
		force = false
	}

	if !force {
		_, err := conf.ReadConfig()
		if err == nil {
			log.Fatal("Configuration already exists. ")
		}
	}

	c := conf.Default

	fmt.Println("Please enter your client id")
	var clientID string
	if _, err := fmt.Scan(&clientID); err != nil {
		log.Fatalf("Unable to read client id %v", err)
	}
	c.Drive.OauthConfig.ClientID = clientID

	fmt.Println("Please enter your client secret")
	var clientSecret string
	if _, err := fmt.Scan(&clientSecret); err != nil {
		log.Fatalf("Unable to read clientSecret %v", err)
	}
	c.Drive.OauthConfig.ClientSecret = clientSecret

	fmt.Println("attempting connection...")
	// Try and make connection and get and save token
	conn, err := dri.GetClient(&c)
	if err != nil {
		log.Fatalln("Could not establish connection to google drive")
	}

	// Ask for upload path in google drive.
	fmt.Println("Please enter the folder in google drive where you keep the backups")
	var folderName string
	if _, err := fmt.Scan(&folderName); err != nil {
		log.Fatalf("Unable to read folder name %v", err)
	}

	folderID, err := dri.FindOrCreateFolder(conn, folderName)
	if err != nil {
		log.Fatalln("Could not find or create folder in google drive", err)
	}

	c.Drive.FolderID = folderID

	fmt.Println("Please enter the local folder to store the backups")
	var localFolder string
	if _, err := fmt.Scan(&localFolder); err != nil {
		log.Fatalf("Unable to read local folder name %v", err)
	}

	c.LocalFolder = localFolder

	if err := conf.WriteConfig(&c); err != nil {
		log.Fatalln("Could not write config file ", err)
	}

}

/*
Description prints a short help message
*/
func (i Init) Description(w io.Writer) {
	fmt.Fprintln(w, "Setup basic configuration")
}
