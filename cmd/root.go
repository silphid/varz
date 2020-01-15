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
  "github.com/silphid/varz/cmd/list"
  "github.com/silphid/varz/common"
  "github.com/spf13/cobra"
  "log"
  "os"
  "path/filepath"

  "github.com/mitchellh/go-homedir"
  "github.com/silphid/varz/cmd/export"
  "github.com/silphid/varz/cmd/get"
  "github.com/silphid/varz/cmd/set"
  "github.com/spf13/viper"
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
  cmd.PersistentFlags().StringVar(&common.Options.ConfigFile, "config", "", "config file (default is $HOME/.varz/config.yaml)")
  cmd.PersistentFlags().StringVar(&common.Options.EnvFile, "env", "", "environment variables file (default is $HOME/.varz/varz.yaml)")

  cmd.AddCommand(export.Cmd)
  cmd.AddCommand(list.Cmd)
  cmd.AddCommand(get.Cmd)
  cmd.AddCommand(set.Cmd)
}

func initConfig() {
  opt := &common.Options
  opt.ConfigDir = getConfigDir(opt.ConfigDir)
  opt.ConfigFile = getConfigFile(opt.ConfigDir, opt.ConfigFile)
  opt.EnvFile = getVarzFile(opt.ConfigDir, opt.EnvFile)
  viper.SetConfigFile(opt.ConfigFile)
  viper.AutomaticEnv()
  if err := viper.ReadInConfig(); err != nil {
    log.Fatalf("failed to read config file %q: %v", opt.ConfigFile, err)
  }
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

  //_, err = fmt.Fprintf(os.Stderr, "Using config dir: %s\n", configDir)
  //if err != nil {
  //  log.Fatalf("failed to output config dir: %s", err)
  //}

  return configDir
}

func getConfigFile(configDir, file string) string {
  if file != "" {
    return file
  }
  file = filepath.Join(configDir, "config.yaml")
  if err := createConfigFileIfNotExist(file); err != nil {
    log.Fatalf("failed to create empty config file %q: %v", file, err)
  }
  return file
}

func getVarzFile(configDir, file string) string {
  if file != "" {
    return file
  }
  file = filepath.Join(configDir, "varz.yaml")
  if err := createVarzFileIfNotExist(file); err != nil {
    log.Fatalf("failed to create empty varz.yaml file %q: %v", file, err)
  }
  return file
}

func createConfigFileIfNotExist(path string) error {
  _, err := os.Stat(path)
  if !os.IsNotExist(err) {
    return nil
  }
  file, err := os.Create(path)
  if err != nil {
    return fmt.Errorf("failed to create empty config path %q: %v", path, err)
  }
  _ = file.Close()
  return nil
}

func createVarzFileIfNotExist(path string) error {
  _, err := os.Stat(path)
  if !os.IsNotExist(err) {
    return nil
  }
  file, err := os.Create(path)
  if err != nil {
    return fmt.Errorf("failed to create empty config path %q: %v", path, err)
  }
  _, err = fmt.Fprint(file,
    `section1:
  ENV_VAR1: "abc"
  ENV_VAR2: 123
  subSection1:
    ENV_VAR3: "cba"
    ENV_VAR4: 321
  subSection2:
    ENV_VAR6: "aaa"

section2:
  ENV_VAR1: "def"
  ENV_VAR2: 456
  subSection1:
    ENV_VAR3: "fed"
    ENV_VAR4: 654
  subSection2:
    ENV_VAR6: "ddd"
`)
  if err != nil {
    return err
  }
  _ = file.Close()
  return nil
}
