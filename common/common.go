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

type mapping = map[interface{}]interface{}

func GetDataFilePath() string {
	return "varz.yaml"
}

func GetVariables(fileName, path string) ([]string, map[string]string, error) {
	section, e := loadSection(fileName, path)
	if e != nil {
		return nil, nil, e
	}

	// Extract and sort valid environment variables
	envVarRegex := regexp.MustCompile(`^[A-Z0-9_]+$`)
	length := len(section)
	names := make([]string, 0, length)
	values := make(map[string]string, length)
	for key := range section {
		keyStr := key.(string)
		if envVarRegex.MatchString(keyStr) {
			names = append(names, keyStr)
			values[keyStr] = fmt.Sprintf("%v", section[keyStr])
		}
	}
	sort.Strings(names)

	return names, values, nil
}

func EnsureSectionExists(filePath, keyPath string) error {
	_, err := loadSection(filePath, keyPath)
	return err
}

func loadSection(fileName string, path string) (mapping, error) {
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
	// Find section within yaml tree
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
			return nil, fmt.Errorf("section not found: %s", path)
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

func GetKeyPathOrDefault(args []string, index int) (string, error) {
	if index < len(args) {
		return args[index], nil
	}
	if !viper.IsSet(defaultKey) {
		return "", errors.New("No key path argument specified and no default set")
	}
	return viper.GetString(defaultKey), nil
}
