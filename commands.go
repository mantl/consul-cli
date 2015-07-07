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
		"health-checks": func() (cli.Command, error) {
			return &command.HealthChecksCommand{
				Meta:	meta,
			}, nil
		},
		"health-node": func() (cli.Command, error) {
			return &command.HealthNodeCommand{
				Meta:	meta,
			}, nil
		},
		"health-service": func() (cli.Command, error) {
			return &command.HealthServiceCommand{
				Meta:	meta,
			}, nil
		},
		"health-state": func() (cli.Command, error) {
			return &command.HealthStateCommand{
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
		"session-create": func() (cli.Command, error) {
			return &command.SessionCreateCommand{
				Meta:	meta,
			}, nil
		},
		"session-destroy": func() (cli.Command, error) {
			return &command.SessionDestroyCommand{
				Meta:	meta,
			}, nil
		},
		"session-info": func() (cli.Command, error) {
			return &command.SessionInfoCommand{
				Meta:	meta,
			}, nil
		},
		"session-list": func() (cli.Command, error) {
			return &command.SessionListCommand{
				Meta:	meta,
			}, nil
		},
		"session-node": func() (cli.Command, error) {
			return &command.SessionNodeCommand{
				Meta:	meta,
			}, nil
		},
		"session-renew": func() (cli.Command, error) {
			return &command.SessionRenewCommand{
				Meta:	meta,
			}, nil
		},
		"status-leader": func() (cli.Command, error) {
			return &command.StatusLeaderCommand{
				Meta:	meta,
			}, nil
		},
		"status-peers": func() (cli.Command, error) {
			return &command.StatusPeersCommand{
				Meta:	meta,
			}, nil
		},
	}
}
