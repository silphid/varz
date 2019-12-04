package list

import (
	"github.com/silphid/varz/cmd/list/env"
	"github.com/silphid/varz/cmd/list/sections"
	"github.com/silphid/varz/cmd/list/entries"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "list sections|entries|env",
	Short: "Lists sections, variables defined in sections, or current environment variables.",
	Long: `Lists sections, variables defined in sections, or current environment variables.`,
}

func init() {
	Cmd.AddCommand(sections.Cmd)
	Cmd.AddCommand(entries.Cmd)
	Cmd.AddCommand(env.Cmd)
}
