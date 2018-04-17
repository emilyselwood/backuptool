package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/wselwood/backuptool/conf"
)

/*
Run creates a backup and stores it.
*/
type Run struct {
}

/*
Description provides help text
*/
func (r Run) Description(w io.Writer) {
	fmt.Fprintln(w, "Create a backup")
}

/*
Run creates a backup according to the current config
*/
func (r Run) Run() {
	c, err := conf.ReadConfig()
	if err != nil {
		log.Fatalln("could not load config ", err)
	}

	log.Println(createTotalFileName())

	f, err := os.Create(createTotalFileName())
	if err != nil {
		log.Fatalln("Could not create output file", err)
	}

	defer f.Close()
	w := zip.NewWriter(f)

	if err := addConfig(w, c); err != nil {
		log.Fatalln("Could not create backup.conf in zip", err)
	}

	// TODO: Thread this
	for _, d := range c.Dirs {
		outName := folderZipName(d.Path)
		log.Println(outName)
		buf, err := zipDir(d)
		if err != nil {
			log.Fatalln("Could not zip dir", d.Path, err)
		}

		f, err := w.Create(outName)
		if err != nil {
			log.Fatalln("could not create file ", outName, err)
		}

		_, err = io.Copy(f, buf)
		if err != nil {
			log.Fatalln("Could not write file", outName, err)
		}
	}
	if err := w.Close(); err != nil {
		log.Fatalln("Could not close zip file", err)
	}
}

func createTotalFileName() string {
	ts := time.Now().Format("2006-01-02T15-04-05")
	return fmt.Sprintf("backup_%s.zip", ts)
}

func folderZipName(d string) string {
	return path.Base(d) + ".zip"
}

func zipDir(d conf.Dir) (*bytes.Buffer, error) {

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	addFiles(w, d, d.Path, "")

	err := w.Close()
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func addFiles(w *zip.Writer, d conf.Dir, basePath, baseInZip string) error {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		p := path.Join(basePath, file.Name())
		if !file.IsDir() {
			if d.ShouldInclude(p) {
				if err := addFile(w, file, basePath, baseInZip); err != nil {
					return err
				}
			}
		} else if file.IsDir() {

			// Recurse
			addFiles(w, d, p, file.Name()+"/")
		}
	}
	return nil
}

func addFile(w *zip.Writer, file os.FileInfo, basePath string, baseInZip string) error {
	dat, err := os.Open(path.Join(basePath, file.Name()))
	if err != nil {
		return err
	}
	defer dat.Close()

	// Add some files to the archive.
	f, err := w.Create(baseInZip + file.Name())
	if err != nil {
		return err
	}
	_, err = io.Copy(f, dat)
	if err != nil {
		return err
	}
	return nil
}

func addConfig(w *zip.Writer, c *conf.Config) error {
	cf, err := w.Create("backup.conf")
	if err != nil {
		return err
	}
	conf.WriteConfigZip(c, cf)
	return nil
}
