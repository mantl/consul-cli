package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/CiscoCloud/consul-cli/commands"
)

const Name = "consul-cli"
const Version = "0.2.0"

func main() {
	log.SetOutput(ioutil.Discard)

	root := commands.Init()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(root.Err, err)
		os.Exit(1)
	}

	os.Exit(0)
	//	args := os.Args[1:]
	//	for _, arg := range args {
	//		if arg == "--" {
	//			break
	//		}
	//
	//		if arg == "-v" || arg == "--version" {
	//			fmt.Printf("%s v%s\n", Name, Version)
	//			os.Exit(0)
	//		}
	//	}
	//
	//	cli := &cli.CLI{
	//		Args:		args,
	//		Commands:	Commands,
	//		HelpFunc:	cli.BasicHelpFunc(Name),
	//	}
	//
	//	exitCode, err := cli.Run()
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
	//		os.Exit(1)
	//	}
	//
	//	os.Exit(exitCode)
}
