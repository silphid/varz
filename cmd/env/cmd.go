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
	Short: "Lists values currently set in your shell environment variables (only those with names corresponding to those in any arbitrary section)",
	Long: `TODO`,
	RunE: run,
	Args: cobra.NoArgs,
}

func run(_ *cobra.Command, _ []string) error {
	names, _, err := common.GetVariables(common.Options.EnvFile, common.ArbitrarySection)
	if err != nil {
		return err
	}

	// Output environment variables
	for _, name := range names {
		value := os.Getenv(name)
		line := fmt.Sprintf("%s=%v\n", name, value)
		fmt.Print(line)
	}

	return nil
}
