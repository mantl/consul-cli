package action

import (
	"flag"
	"fmt"
)

type healthNode struct {
	*config
}

func HealthNodeAction() Action {
	return &healthNode{
		config: &gConfig,
	}
}

func (h *healthNode) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	h.addDatacenterFlag(f)
	h.addOutputFlags(f, false)
	h.addConsistencyFlags(f)

	return f
}

func (h *healthNode) Run(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one node name allowed")
	}
	node := args[0]

	client, err := h.newHealth()
	if err != nil {
		return err
	}

	queryOpts := h.queryOptions()

	n, _, err := client.Node(node, queryOpts)
	if err != nil {
		return err
	}

	return h.Output(n)
}
