package action

import (
	"flag"
)

type Action interface {
	CommandFlags() *flag.FlagSet
	Run([]string) error
}

func newFlagSet() *flag.FlagSet {
	return flag.NewFlagSet("consul-cli", flag.ExitOnError)
}

func GlobalCommandFlags() *flag.FlagSet {
	return gFlags
}
