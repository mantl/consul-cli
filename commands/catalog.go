package commands

import "github.com/spf13/cobra"

type Catalog struct {
	*Cmd
}

func (root *Cmd) initCatalog() {
	c := Catalog{Cmd: root}

	catalogCmd := &cobra.Command{
		Use:   "catalog",
		Short: "Consul /catalog endpoint interface",
		Long:  "Consul /catalog endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	c.AddDatacentersSub(catalogCmd)
	c.AddNodeSub(catalogCmd)
	c.AddNodesSub(catalogCmd)
	c.AddServiceSub(catalogCmd)
	c.AddServicesSub(catalogCmd)
	c.AddDeregisterSub(catalogCmd)

	c.AddCommand(catalogCmd)
}
