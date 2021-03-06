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
  "github.com/mitchellh/go-homedir"
  "github.com/silphid/varz/cmd/entries"
  "github.com/silphid/varz/cmd/env"
  "github.com/silphid/varz/cmd/export"
  "github.com/silphid/varz/cmd/sections"
  "github.com/silphid/varz/common"
  "github.com/spf13/cobra"
  "log"
  "os"
  "path/filepath"
)

var cmd = &cobra.Command{
  Use:   "varz",
  Short: "Allows to quickly export different sets of environment variables to current shell",
  Long: `Varz allows to quickly export different sets of environment variables to current shell.`,
}

func Execute() {
  if err := cmd.Execute(); err != nil {
      os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  cmd.PersistentFlags().StringVar(&common.Options.ConfigDir, "config-dir", "", "config dir (default is $HOME/.varz)")
  cmd.PersistentFlags().StringVar(&common.Options.EnvFile, "env", "", "environment variables file (default is $HOME/.varz/varz.yaml)")

  cmd.AddCommand(export.Cmd)
  cmd.AddCommand(sections.Cmd)
  cmd.AddCommand(entries.Cmd)
  cmd.AddCommand(env.Cmd)
}

func initConfig() {
  opt := &common.Options
  opt.ConfigDir = getConfigDir(opt.ConfigDir)
  opt.EnvFile = getVarzFile(opt.ConfigDir, opt.EnvFile)
}

func getConfigDir(configDir string) string {
  if configDir != "" {
    return configDir
  }

  // VARZ env variable defines directory to use
  configDir, ok := os.LookupEnv("VARZ")
  if !ok {
    // Fallback to "~/.varz"
    homeDir, err := homedir.Dir()
    if err != nil {
      log.Fatalf("failed to find $HOME directory: %v", err)
    }
    configDir = filepath.Join(homeDir, ".varz")
  }

  // Resolve to absolute path
  configDir, err := filepath.Abs(configDir)
  if err != nil {
    log.Fatalf("failed to resolve absolute path of config dir %s: %v", configDir, err)
  }

  // Ensure directory exists or create it
  if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
    log.Fatalf("failed to create config dir %q: %v", configDir, err)
  }

  return configDir
}

func getVarzFile(configDir, file string) string {
  if file != "" {
    return file
  }
  file = filepath.Join(configDir, "varz.yaml")
  _, err := os.Stat(file)
  if os.IsNotExist(err) {
    log.Fatalf("missing varz.yaml file: %v", err)
  }
  return file
}
