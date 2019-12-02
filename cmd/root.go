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
  "github.com/spf13/cobra"
  "os"
  "path/filepath"

  homedir "github.com/mitchellh/go-homedir"
  "github.com/silphid/varz/cmd/dump"
  "github.com/silphid/varz/cmd/export"
  "github.com/silphid/varz/cmd/get"
  "github.com/silphid/varz/cmd/list"
  "github.com/silphid/varz/cmd/set"
  "github.com/spf13/viper"
)

var cfgFile string

// cmd represents the base command when called without any subcommands
var cmd = &cobra.Command{
  Use:   "varz",
  Short: "Allows to quickly export different sets of environment variables to current shell",
  Long: `TODO...`,
  // Uncomment the following line if your bare application
  // has an action associated with it:
  //	Run: func(cmd *cobra.Command, args []string) { },
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

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.

  cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (_default is $HOME/.varz.yaml)")


  // Cobra also supports local flags, which will only run
  // when this action is called directly.
  cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

  cmd.AddCommand(export.Cmd)
  cmd.AddCommand(list.Cmd)
  cmd.AddCommand(dump.Cmd)
  cmd.AddCommand(get.Cmd)
  cmd.AddCommand(set.Cmd)
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      _, _ = fmt.Fprintf(os.Stderr, "Failed to find home directory: %v", err)
      os.Exit(1)
    }

    configFile := filepath.Join(home, ".varz.yaml")
    viper.SetConfigFile(configFile)

    // Force creation of config file, if not already existing
    if !fileExists(configFile) {
      emptyFile, err := os.Create(configFile)
      if err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "Failed to create config file %s: %v", configFile, err)
        os.Exit(1)
      }
      _ = emptyFile.Close()
    }
  }

  viper.AutomaticEnv() // read in environment variables that match
  _ = viper.ReadInConfig() // if a config file is found, read it in.
}

func fileExists(filename string) bool {
  info, err := os.Stat(filename)
  if os.IsNotExist(err) {
    return false
  }
  return !info.IsDir()
}