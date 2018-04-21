package main

import (
	"fmt"
	"io"
	"log"

	"parsecsreach.com/backuptool/conf"
)

/*
List command
*/
type List struct {
}

/*
Description provides help text
*/
func (l List) Description(w io.Writer) {
	fmt.Fprintln(w, "List the directories to be backed up")
}

/*
Run Lists the current backup config
*/
func (l List) Run() {
	c, err := conf.ReadConfig()
	if err != nil {
		log.Fatalln("could not load config", err)
	}

	if len(c.Dirs) > 0 {
		fmt.Println("Path\tInclude\tExclude")
	}

	for _, d := range c.Dirs {
		fmt.Println(d)
	}
}
