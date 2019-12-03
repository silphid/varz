package _default

import (
	"github.com/silphid/varz/common"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "default",
	Short: "Sets default key path to use",
	Long: `TODO...`,
	RunE: run,
	Args: cobra.ExactArgs(1),
}

func run(_ *cobra.Command, args []string) error {
	keyPath := args[0]
	if err := common.EnsureSectionExists(common.Options.DataFile, keyPath); err != nil {
		return err
	}
	return common.SetDefaultKeyPath(keyPath)
}