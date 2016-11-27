package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ChrisAubuchon/consul-cli/commands"
)

const Name = "consul-cli"
const Version = "0.3.1"

func main() {
	log.SetOutput(ioutil.Discard)

	root := commands.Init(Name, Version)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(root.Err, err)
		os.Exit(1)
	}

	os.Exit(0)
}
