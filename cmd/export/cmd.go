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
package export

import (
	"fmt"
	"os"
	"strings"

	"github.com/silphid/varz/common"
	"github.com/spf13/cobra"
)

var verbose *bool

func init() {
	verbose = Cmd.Flags().BoolP("verbose", "v", false, "Outputs more information on stderr")
}

var Cmd = &cobra.Command {
	Use:   "export [SECTION]",
	Short: "Outputs export statements for given subset of variables",
	Long: `The output of this command is intended to be sourced
in order to define the corresponding environment variables
in your current shell. For example:

  . <(varz export SECTION)`,
	RunE: run,
	Args: cobra.RangeArgs(0, 1),
}

func run(_ *cobra.Command, args []string) error {
	keyPath := ""
	if len(args) == 1 {
		keyPath = args[0]
	}

	stdout, stderr, err := do(common.Options.EnvFile, keyPath)
	if err != nil {
		return err
	}

	if stdout != "" {
		fmt.Print(stdout)
	}

	if stderr != "" {
		if _, err := fmt.Fprint(os.Stderr, stderr); err != nil {
			return err
		}
	}

	return nil
}

func do(file, keyPath string) (string, string, error) {
	keyPath, err := common.GetKeyPathOrDefault(keyPath)
	if err != nil {
		return "", "", err
	}
	names, values, err := common.GetVariables(file, keyPath)
	if err != nil {
		return "", "", err
	}

	// Output variables
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	for _, name := range names {
		line := fmt.Sprintf("export %s=%v\n", name, values[name])
		stdout.WriteString(line)
		if *verbose {
			stderr.WriteString(line)
		}
	}

	return stdout.String(), stderr.String(), nil
}
