package common

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

type GlobalOptions struct {
	ConfigDir  string
	ConfigFile string
	EnvFile    string
}

var Options GlobalOptions

type mapping = map[interface{}]interface{}

func GetVariables(fileName, path string) ([]string, map[string]string, error) {
	section, e := LoadSection(fileName, path)
	if e != nil {
		return nil, nil, e
	}

	// Extract and sort valid environment variables
	length := len(section)
	names := make([]string, 0, length)
	values := make(map[string]string, length)
	for key := range section {
		keyStr := key.(string)
		if IsEnvVarName(keyStr) {
			names = append(names, keyStr)
			values[keyStr] = fmt.Sprintf("%v", section[keyStr])
		}
	}
	sort.Strings(names)

	return names, values, nil
}

func GetSections(fileName, path string) ([]string, error) {
	section, e := LoadSection(fileName, path)
	if e != nil {
		return nil, e
	}

	// Extract and sort valid sections names
	length := len(section)
	names := make([]string, 0, length)
	for key := range section {
		keyStr := key.(string)
		if !IsEnvVarName(keyStr) {
			names = append(names, keyStr)
		}
	}

	return names, nil
}

var envVarRegex = regexp.MustCompile(`^[A-Z0-9_]+$`)

func IsEnvVarName(key string) bool {
	return envVarRegex.MatchString(key)
}

func EnsureSectionExists(filePath, keyPath string) error {
	_, err := LoadSection(filePath, keyPath)
	return err
}

func LoadSection(fileName string, path string) (mapping, error) {
	// Load file to buffer
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// Parse buffer as yaml into map
	m := make(mapping)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	// Find sections within yaml tree
	if path == "" {
		return m, nil
	}
	section, err := getSection(m, path)
	if err != nil {
		return nil, err
	}
	return section, nil
}

func getSection(m mapping, path string) (mapping, error) {
	cur := m
	for _, component := range strings.Split(path, "/") {
		val, ok := cur[component].(mapping)
		if !ok {
			return nil, fmt.Errorf("sections not found: %s", path)
		}
		cur = val
	}
	return cur, nil
}

var defaultKey = "default"

func SetDefaultKeyPath(value string) error {
	viper.Set(defaultKey, value)
	if err := viper.WriteConfig(); err != nil {
		return errors.Wrap(err, "failed to write default key path")
	}
	return nil
}

func GetDefaultKeyPath() (string, error) {
	if !viper.IsSet(defaultKey) {
		return "", errors.New("No default key path already set")
	}
	return viper.GetString(defaultKey), nil
}

func GetKeyPathOrDefault(keyPath string) (string, error) {
	if keyPath != "" {
		return keyPath, nil
	}
	if !viper.IsSet(defaultKey) {
		return "", errors.New("No key path argument specified and no default set")
	}
	return viper.GetString(defaultKey), nil
}
