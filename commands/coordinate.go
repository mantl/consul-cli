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
		Short: "Consul Coordinate endpoint interface",
		Long:  "Consul Coordinate endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			root.Help()
		},
	}

	c.AddDatacentersSub(coordinateCmd)
	c.AddNodesSub(coordinateCmd)

	c.AddCommand(coordinateCmd)
}
