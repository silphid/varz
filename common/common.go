package common

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

type mapping = map[interface{}]interface{}

func GetVariables(fileName, path string) ([]string, map[string]string, error) {
	// Load file to buffer
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, nil, err
	}

	// Parse buffer as yaml into map
	m := make(mapping)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, nil, err
	}

	// Find section within yaml tree
	section, err := getSection(m, path)
	if err != nil {
		return nil, nil, err
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
			values[keyStr] = section[keyStr].(fmt.Stringer).String()
		}
	}
	sort.Strings(names)

	return names, values, nil
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