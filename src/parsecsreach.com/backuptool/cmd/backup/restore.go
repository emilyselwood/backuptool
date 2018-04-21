package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"parsecsreach.com/backuptool/conf"
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
	config, err := conf.ReadConfigZip(rootFile)
	if err != nil {
		log.Fatalln("Could not extract config", err)
	}

	for _, d := range rootFile.File {
		if d.Name != "backup.conf" {
			log.Println("extracting ", d.Name)
			// Find matching config
			dir := findMatching(config, d.Name)
			// Create directories
			// TODO: file perms need thinking about.
			basePath := path.Join(prefix, dir.Path)
			os.MkdirAll(basePath, os.ModeExclusive)

			// extract files into new directories
			if err := extractZip(d, basePath, overwrite); err != nil {
				log.Fatalln("Could not extract into ", basePath, err)
			}
		}

	}

}

func findMatching(config *conf.Config, name string) *conf.Dir {
	n := strings.TrimSuffix(name, ".zip")
	for _, d := range config.Dirs {
		if strings.HasSuffix(d.Path, n) {
			return &d
		}
	}
	return nil
}

func extractZip(parent *zip.File, basePath string, overwrite bool) error {

	// To read a zip file nested inside another zip file first we have to uncompress the inner zip file
	// as the zip.NewReader needs a ReadSeaker rather than a readCloser we get from parent.Open()
	f, err := parent.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	var buffer bytes.Buffer
	io.Copy(&buffer, f)

	// Here we have to convert the buffer into a reader so that we get the ReadSeaker interface.
	// We might have been able to get around this by using something with stream compression (.gz) but
	// Zip is easier to work with on more machines, if a single file needs to be extracted from a backup
	reader, err := zip.NewReader(bytes.NewReader(buffer.Bytes()), int64(buffer.Len()))
	if err != nil {
		return err
	}

	for _, r := range reader.File {
		if r.FileInfo().IsDir() {
			if err := saveDir(r, basePath); err != nil {
				return err
			}
		} else {
			if err := saveFile(r, basePath, overwrite); err != nil {
				return err
			}
		}

	}
	return nil
}

func saveDir(r *zip.File, basePath string) error {
	return os.MkdirAll(path.Join(basePath, r.Name), r.Mode())
}

func saveFile(r *zip.File, basePath string, overwrite bool) error {
	p := filepath.Join(basePath, r.Name)
	// Do we need to defer files until the directory is done so the dir permissions don't get messed up?
	os.MkdirAll(filepath.Dir(p), r.Mode())

	mode := os.O_CREATE
	if !overwrite {
		mode = mode | os.O_EXCL
	}
	n, err := os.OpenFile(p, mode, r.Mode())
	if err != nil {
		if !overwrite && os.IsExist(err) {
			return nil
		}
		return err
	}
	defer n.Close()

	i, err := r.Open()
	if err != nil {
		return err
	}
	defer i.Close()

	_, err = io.Copy(n, i)

	return err
}
