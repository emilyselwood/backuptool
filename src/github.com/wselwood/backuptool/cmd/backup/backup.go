package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

/*
Cmd is the type for a command that can be launched from the command line
*/
type Cmd interface {
	Run()
	Description(w io.Writer)
}

var commands = map[string]Cmd{
	"help":    Help{},
	"init":    Init{},
	"create":  Run{},
	"add":     Add{},
	"list":    List{},
	"restore": Restore{},
}

func main() {
	key := ""
	if len(os.Args) > 1 {
		key = strings.ToLower(os.Args[1])
	}

	c, ok := commands[key]
	if !ok {
		log.Println("Could not understand command ", key)
		commands["help"].Run()
	} else {
		c.Run()
	}

}

/*
Help is a command that prints out the help text
*/
type Help struct {
}

/*
Run method for help prints out the help text
*/
func (h Help) Run() {
	log.Println("help")
	for k, v := range commands {
		fmt.Printf("%s: ", k)
		v.Description(os.Stdout)
	}
}

/*
Description prints a short help message
*/
func (h Help) Description(w io.Writer) {
	fmt.Fprintln(w, "help command")
}
