package main

import (
	"os"

	"github.com/CiscoCloud/consul-cli/command"
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func init() {
	metaPtr := new(command.Meta)
	meta := *metaPtr
	meta.UI = &cli.BasicUi{
		Writer:		os.Stdout,
		ErrorWriter:	os.Stderr,
	}

	Commands = map[string]cli.CommandFactory{
		"acl-clone": func() (cli.Command, error) {
			return &command.ACLCloneCommand{
				Meta:	meta,
			}, nil
		},
		"acl-create": func() (cli.Command, error) {
			return &command.ACLCreateCommand{
				Meta:	meta,
			}, nil
		},
		"acl-destroy": func() (cli.Command, error) {
			return &command.ACLDestroyCommand{
				Meta:	meta,
			}, nil
		},
		"acl-info": func() (cli.Command, error) {
			return &command.ACLInfoCommand{
				Meta:	meta,
			}, nil
		},
		"acl-list": func() (cli.Command, error) {
			return &command.ACLListCommand{
				Meta:	meta,
			}, nil
		},
		"acl-update": func() (cli.Command, error) {
			return &command.ACLUpdateCommand{
				Meta:	meta,
			}, nil
		},
		"kv-delete": func() (cli.Command, error) {
			return &command.KVDeleteCommand{
				Meta:	meta,
			}, nil
		},

		"kv-read": func() (cli.Command, error) {
			return &command.KVReadCommand{
				Meta:	meta,
			}, nil
		},

		"kv-write": func() (cli.Command, error) {
			return &command.KVWriteCommand{
				Meta:	meta,
			}, nil
		},

		"kv-lock": func() (cli.Command, error) {
			return &command.KVLockCommand{
				Meta:	meta,
			}, nil
		},

		"kv-unlock": func() (cli.Command, error) {
			return &command.KVUnlockCommand{
				Meta:	meta,
			}, nil
		},
		"session-list": func() (cli.Command, error) {
			return &command.SessionListCommand{
				Meta:	meta,
			}, nil
		},
	}
}
