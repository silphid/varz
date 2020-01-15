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

const BaseKey string = "base"

var Options GlobalOptions

type Document = map[interface{}]interface{}

func GetVariables(fileName, path string) ([]string, map[string]string, error) {
	doc, err := Load(fileName)
	if err != nil {
		return nil, nil, err
	}

	baseSection, _ := GetSection(doc, BaseKey)
	requestedSection, err := GetSection(doc, path)
	if err != nil {
		return nil, nil, err
	}

	// Extract and sort valid environment variables
	length := len(baseSection) + len(requestedSection)
	names := make([]string, 0, length)
	values := make(map[string]string, length)
	for _, section := range []Document{baseSection, requestedSection} {
		for key := range section {
			keyStr := key.(string)
			if IsEnvVarName(keyStr) {
				names = append(names, keyStr)
				values[keyStr] = fmt.Sprintf("%v", section[keyStr])
			}
		}
	}
	sort.Strings(names)

	return names, values, nil
}

func GetSections(fileName, path string) ([]string, error) {
	doc, e := Load(fileName)
	if e != nil {
		return nil, e
	}

	section, err := GetSection(doc, path)
	if err != nil {
		return nil, err
	}

	// Extract and sort valid sections names
	length := len(section)
	names := make([]string, 0, length)
	for key := range section {
		keyStr := key.(string)
		if !IsEnvVarName(keyStr) && keyStr != BaseKey {
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
	doc, err := Load(filePath)
	if err != nil {
		return err
	}

	_, err = GetSection(doc, keyPath)
	return err
}

func Load(fileName string) (Document, error) {
	// Load file to buffer
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// Parse buffer as yaml into map
	doc := make(Document)
	err = yaml.Unmarshal(data, &doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func GetSection(doc Document, path string) (Document, error) {
	if path == "" {
		return doc, nil
	}

	cur := doc
	for _, component := range strings.Split(path, "/") {
		val, ok := cur[component].(Document)
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

func GetKeyPathOrDefault(keyPath string) (string, error) {
	if keyPath != "" {
		return keyPath, nil
	}
	if !viper.IsSet(defaultKey) {
		return "", errors.New("No key path argument specified and no default set")
	}
	return viper.GetString(defaultKey), nil
}
