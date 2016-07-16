package commands

import (
	"github.com/spf13/cobra"
)

type Coordinate struct {
	*Cmd
}

func (root *Cmd) initCoordinate() {
	c := Coordinate{Cmd: root}

	coordinateCmd := &cobra.Command{
		Use:   "coordinate",
		Short: "Consul /coordinate endpoint interface",
		Long:  "Consul /coordinate endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	c.AddDatacentersSub(coordinateCmd)
	c.AddNodesSub(coordinateCmd)

	c.AddCommand(coordinateCmd)
}
