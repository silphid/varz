package show

import (
	"github.com/silphid/varz/cmd/show/env"
	"github.com/silphid/varz/cmd/show/file"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "show",
	Short: "show env|file",
	Long: `TODO...`,
}

func init() {
	Cmd.AddCommand(env.Cmd)
	Cmd.AddCommand(file.Cmd)
}
