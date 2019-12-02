package set

import (
	"github.com/silphid/varz/cmd/set/_default"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "set",
	Short: "set default",
	Long: `TODO...`,
}

func init() {
	Cmd.AddCommand(_default.Cmd)
}
