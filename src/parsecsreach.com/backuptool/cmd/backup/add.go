package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"parsecsreach.com/backuptool/conf"
)

/*
Add creates a backup and stores it.
*/
type Add struct {
}

/*
Description provides help text
*/
func (a Add) Description(w io.Writer) {
	fmt.Fprintln(w, "Add a directory to be backed up")
}

/*
Run creates a backup according to the current config
*/
func (a Add) Run() {

	var args []string
	var inc string
	var exc string

	if len(os.Args) > 3 {
		mySet := flag.NewFlagSet("", flag.ExitOnError)
		mySet.StringVar(&inc, "i", "", "include regex")
		mySet.StringVar(&exc, "e", "", "exclude regex")
		mySet.Parse(os.Args[2:])
		args = mySet.Args()
	} else {
		args = os.Args[2:]
		inc = ""
		exc = ""
	}

	newDirs := []conf.Dir{}

	for _, d := range args {
		n := conf.Dir{
			Path:    d,
			Include: inc,
			Exclude: exc,
		}
		log.Println(n)
		newDirs = append(newDirs, n)

	}

	c, err := conf.ReadConfig()
	if err != nil {
		log.Fatalln("could not load config", err)
	}

	if c.Dirs == nil {
		c.Dirs = newDirs
	} else {
		c.Dirs = append(c.Dirs, newDirs...)
	}

	if err := conf.WriteConfig(c); err != nil {
		log.Fatalln("Could not write config", err)
	}

}
