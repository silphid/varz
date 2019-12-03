/*
Copyright © 2019 Mathieu Frenette <mathieu@silphid.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package env

import (
	"fmt"
	"os"

	"github.com/silphid/varz/common"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command {
	Use:   "env",
	Short: "Shows current values in your shell environment for the given variables set",
	Long: `TODO`,
	RunE: run,
	Args: cobra.RangeArgs(0, 1),
}

func run(_ *cobra.Command, args []string) error {
	keyPath := ""
	if len(args) == 1 {
		keyPath = args[0]
	}
	keyPath, err := common.GetKeyPathOrDefault(keyPath)
	if err != nil {
		return err
	}
	names, _, err := common.GetVariables(common.Options.EnvFile, keyPath)
	if err != nil {
		return err
	}

	// Output environment variables
	for _, name := range names {
		value := os.Getenv(name)
		line := fmt.Sprintf("%s=%v\n", name, value)
		fmt.Printf(line)
	}

	return nil
}
