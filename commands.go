// +build

package main

import (
	"os"

	"github.com/CiscoCloud/consul-cli/commands"
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
		"agent-checks": func() (cli.Command, error) {
			return &command.AgentChecksCommand{
				Meta:	meta,
			}, nil
		},
		"agent-force-leave": func() (cli.Command, error) {
			return &command.AgentForceLeaveCommand{
				Meta:	meta,
			}, nil
		},
		"agent-join": func() (cli.Command, error) {
			return &command.AgentJoinCommand{
				Meta:	meta,
			}, nil
		},
		"agent-maintenance": func() (cli.Command, error) {
			return &command.AgentMaintenanceCommand{
				Meta:	meta,
			}, nil
		},
		"agent-members": func() (cli.Command, error) {
			return &command.AgentMembersCommand{
				Meta:	meta,
			}, nil
		},
		"agent-self": func() (cli.Command, error) {
			return &command.AgentSelfCommand{
				Meta:	meta,
			}, nil
		},
		"agent-services": func() (cli.Command, error) {
			return &command.AgentServicesCommand{
				Meta:	meta,
			}, nil
		},
		"catalog-datacenters": func() (cli.Command, error) {
			return &command.CatalogDatacentersCommand{
				Meta:	meta,
			}, nil
		},
		"catalog-nodes": func() (cli.Command, error) {
			return &command.CatalogNodesCommand{
				Meta:	meta,
			}, nil
		},
		"catalog-node": func() (cli.Command, error) {
			return &command.CatalogNodeCommand{
				Meta:	meta,
			}, nil
		},
		"catalog-services": func() (cli.Command, error) {
			return &command.CatalogServicesCommand{
				Meta:	meta,
			}, nil
		},
		"catalog-service": func() (cli.Command, error) {
			return &command.CatalogServiceCommand{
				Meta:	meta,
			}, nil
		},
		"check-fail": func() (cli.Command, error) {
			return &command.CheckFailCommand{
				Meta:	meta,
			}, nil
		},
		"check-deregister": func() (cli.Command, error) {
			return &command.CheckDeregisterCommand{
				Meta:	meta,
			}, nil
		},
		"check-pass": func() (cli.Command, error) {
			return &command.CheckPassCommand{
				Meta:	meta,
			}, nil
		},
		"check-register": func() (cli.Command, error) {
			return &command.CheckRegisterCommand{
				Meta:	meta,
			}, nil
		},
		"check-warn": func() (cli.Command, error) {
			return &command.CheckWarnCommand{
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
		"kv-watch": func() (cli.Command, error) {
			return &command.KVWatchCommand{
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
		"service-deregister": func() (cli.Command, error) {
			return &command.ServiceDeregisterCommand{
				Meta:	meta,
			}, nil
		},
		"service-maintenance": func() (cli.Command, error) {
			return &command.ServiceMaintenanceCommand{
				Meta:	meta,
			}, nil
		},
		"service-register": func() (cli.Command, error) {
			return &command.ServiceRegisterCommand{
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
