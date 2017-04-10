package action

import (
	"flag"
)

type Action interface {
	CommandFlags() *flag.FlagSet
	Run([]string) error
}

func GlobalCommandFlags() *flag.FlagSet {
	return gFlags
}
