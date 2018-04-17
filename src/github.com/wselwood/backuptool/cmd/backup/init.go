package main

import (
	"fmt"
	"io"
	"log"

	"github.com/wselwood/backuptool/conf"
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

	_, err := conf.ReadConfig()
	if err == nil {
		log.Fatal("Configuration already exists. ")
	}

	c := conf.Default

	// TODO: Ask for keys and so on.

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
