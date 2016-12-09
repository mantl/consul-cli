package commands

import (
	"github.com/spf13/cobra"
)

// getStringSlice works around viper.GetStringSlice's bad behavior.
// viper.GetStringSlice gets the flag with viper.Get() and casts the
// result use spf13.cast.ToStringSlice(). The issue is that viper.Get()
// returns an interface{} which cast.ToStringSlice() converts to
// []string{value}
//
// so, for example, `--arg=foo --arg=bar` using viper.GetStringSlice()
// will return []string{"foo,bar"} which is clearly incorrect.
//
func getStringSlice(cmd *cobra.Command, flag string) []string {
	ss, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		return []string{}
	}

	return ss
}
