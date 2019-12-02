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
package list

import (
	"fmt"
	"varz/common"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command {
	Use:   "list",
	Short: "Outputs export statements for given subset of variables",
	Long: `The output of this command is intended to be sourced
in order to define the corresponding environment variables
in your current shell. For example:

. <(varz export path/to/vars)`,
	RunE: run,
	Args: cobra.ExactArgs(1),
}

func run(_ *cobra.Command, args []string) error {
	names, values, err := common.GetVariables("varz.yaml", args[0])
	if err != nil {
		return err
	}

	// Output environment variables
	for _, name := range names {
		line := fmt.Sprintf("%s=%v\n", name, values[name])
		fmt.Printf(line)
	}

	return nil
}
