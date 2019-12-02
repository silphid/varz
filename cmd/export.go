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
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type mapping = map[interface{}]interface{}

var verbose *bool

func init() {
	rootCmd.AddCommand(exportCmd)
	verbose = exportCmd.Flags().BoolP("verbose", "v", false, "Outputs more information on stderr")
}

var exportCmd = &cobra.Command {
	Use:   "export",
	Short: "Outputs export statements for given subset of variables",
	Long: `The output of this command is intended to be sourced
in order to define the corresponding environment variables
in your current shell. For example:

. <(varz export path/to/vars)`,
	RunE: run,
	Args: cobra.ExactArgs(1),
}

func run(_ *cobra.Command, args []string) error {
	// Load file to buffer
	data, err := ioutil.ReadFile("varz.yaml")
	if err != nil {
		return err
	}

	// Parse buffer as yaml into map
	m := make(mapping)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	// Find section within yaml tree
	section, err := getSection(m, args[0])
	if err != nil {
		return err
	}

	// Extract and sort valid environment variables
	envVarRegex := regexp.MustCompile(`^[A-Z0-9_]+$`)
	keys := make([]string, 0, len(section))
	for key := range section {
		keyStr := key.(string)
		if envVarRegex.MatchString(keyStr) {
			keys = append(keys, keyStr)
		}
	}
	sort.Strings(keys)

	// Output environment variables
	for _, key := range keys {
		line := fmt.Sprintf("export %s=%v\n", key, section[key])
		fmt.Printf(line)
		if *verbose {
			if _, err := fmt.Fprintf(os.Stderr, line); err != nil {
				return err
			}
		}
	}

	return nil
}

func getSection(m mapping, path string) (mapping, error) {
	cur := m
	for _, component := range strings.Split(path, "/") {
		val, ok := cur[component].(mapping)
		if !ok {
			return nil, fmt.Errorf("section not found: %s", path)
		}
		cur = val
	}
	return cur, nil
}