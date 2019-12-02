package get

import (
	"github.com/silphid/varz/cmd/get/_default"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "get default",
	Long: `TODO...`,
}

func init() {
	Cmd.AddCommand(_default.Cmd)
}
