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
  "github.com/silphid/varz/cmd/show/env"
  "github.com/silphid/varz/cmd/show/file"
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

// cmd represents the base command when called without any subcommands
var cmd = &cobra.Command{
  Use:   "varz",
  Short: "Allows to quickly export different sets of environment variables to current shell",
  Long: `TODO...`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the cmd.
func Execute() {
  if err := cmd.Execute(); err != nil {
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  cmd.PersistentFlags().StringVar(&common.Options.ConfigDir, "config-dir", "", "config dir (default is $HOME/.varz)")
  cmd.PersistentFlags().StringVar(&common.Options.ConfigFile, "config", "", "config file (default is $HOME/.varz/config.yaml)")
  cmd.PersistentFlags().StringVar(&common.Options.EnvFile, "env", "", "environment variables file (default is $HOME/.varz/env.varz)")

  cmd.AddCommand(export.Cmd)
  cmd.AddCommand(file.Cmd)
  cmd.AddCommand(env.Cmd)
  cmd.AddCommand(get.Cmd)
  cmd.AddCommand(set.Cmd)
}

func initConfig() {
  opt := &common.Options
  opt.ConfigDir = getConfigDir(opt.ConfigDir)
  opt.ConfigFile = getConfigFile(opt.ConfigDir, opt.ConfigFile)
  opt.EnvFile = getEnvFile(opt.ConfigDir, opt.EnvFile)
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

  homeDir, err := homedir.Dir()
  if err != nil {
    log.Fatalf("failed to find $HOME directory: %v", err)
  }

  configDir = filepath.Join(homeDir, ".varz")
  if err := createDirIfNotExist(configDir); err != nil {
    log.Fatalf("failed to create config dir %q: %v", configDir, err)
  }

  return configDir
}

func getConfigFile(configDir, file string) string {
  if file != "" {
    return file
  }
  file = filepath.Join(configDir, "config.yaml")
  if err := createFileIfNotExist(file); err != nil {
    log.Fatalf("failed to create empty config file %q: %v", file, err)
  }
  return file
}

func getEnvFile(configDir, file string) string {
  if file != "" {
    return file
  }
  file = filepath.Join(configDir, "env.varz")
  if err := createFileIfNotExist(file); err != nil {
    log.Fatalf("failed to create empty env file %q: %v", file, err)
  }
  return file
}

func createDirIfNotExist(path string) error {
  info, err := os.Stat(path)
  if os.IsExist(err) {
    return nil
  }
  if err != nil {
    return err
  }
  if !info.IsDir() {
    return fmt.Errorf("expecting directory, but found file %q: %v", path, err)
  }
  if err := os.Mkdir(path, os.ModeDir | os.ModePerm); err != nil {
    return err
  }
  return nil
}

func createFileIfNotExist(path string) error {
  info, err := os.Stat(path)
  if os.IsExist(err) {
    return nil
  }
  if err != nil {
    return err
  }
  if info.IsDir() {
    return fmt.Errorf("expecting file, but found directory %q: %v", path, err)
  }
  file, err := os.Create(path)
  if err != nil {
    return fmt.Errorf("failed to create empty config path %q: %v", path, err)
  }
  _ = file.Close()
  return nil
}
