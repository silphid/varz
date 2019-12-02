package _default

import (
	"fmt"
	"github.com/silphid/varz/common"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "default",
	Short: "Gets default configured key path",
	Long: `TODO...`,
	RunE: run,
	Args: cobra.ExactArgs(0),
}

func run(_ *cobra.Command, _ []string) error {
	value, err := common.GetDefaultKeyPath()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", value)
	return nil
}