package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ChrisAubuchon/consul-cli/commands"
)

const Name = "consul-cli"
const Version = "0.4.0"

func main() {
	log.SetOutput(ioutil.Discard)

	root := commands.Init(Name, Version)
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
