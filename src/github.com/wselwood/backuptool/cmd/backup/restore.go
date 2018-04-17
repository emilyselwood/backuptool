package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

/*
Restore is a command to restore a backup.
*/
type Restore struct {
}

/*
Description provides help text
*/
func (r Restore) Description(w io.Writer) {
	fmt.Fprintln(w, "List the directories to be backed up")
}

/*
Run restores the backup
*/
func (r Restore) Run() {
	// args for prefix (for testing)
	var args []string
	var prefix string
	var overwrite bool

	if len(os.Args) > 3 {
		mySet := flag.NewFlagSet("", flag.ExitOnError)
		mySet.StringVar(&prefix, "p", "", "prefix")
		mySet.BoolVar(&overwrite, "o", false, "overwrite existing files")
		mySet.Parse(os.Args[2:])
		args = mySet.Args()
	} else {
		args = os.Args[2:]
		prefix = "/"
		overwrite = false
	}

	// open zip file
	rootFile, err := zip.OpenReader(args[0])
	if err != nil {
		log.Fatalln("Could not open file ", args[0], err)
	}

	defer rootFile.Close()

	// Extract config
	/*config, err := conf.ReadConfigZip(rootFile)
	if err != nil {
		log.Fatalln("Could not extract config", err)
	}*/

	for _, d := range rootFile.File {
		if d.Name != "backup.conf" {
			log.Println("extracting ", d.Name)
			// TODO: Find matching config
			// TODO: Create directories
			// TODO: extract files into new directories

		}

	}

}
