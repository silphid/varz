package common

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sort"
)

type GlobalOptions struct {
	ConfigDir  string
	EnvFile    string
}

const BaseKey string = "base"
const ArbitrarySection = "???"

var Options GlobalOptions

type Document map[interface{}]interface{}

func (doc Document) GetSection(name string) (Document, error) {
	if name == ArbitrarySection {
		for key, value := range doc {
			if key != BaseKey {
				val, ok := value.(Document)
				if ok {
					return val, nil
				}
			}
		}
	}

	val, ok := doc[name].(Document)
	if !ok {
		return nil, fmt.Errorf("section not found: %s", name)
	}
	return val, nil
}

func GetVariables(fileName, sectionName string) ([]string, map[string]string, error) {
	doc, err := LoadDocument(fileName)
	if err != nil {
		return nil, nil, err
	}

	baseSection, _ := doc.GetSection(BaseKey)
	requestedSection, err := doc.GetSection(sectionName)
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
			names = append(names, keyStr)
			values[keyStr] = fmt.Sprintf("%v", section[keyStr])
		}
	}
	sort.Strings(names)

	return names, values, nil
}

func LoadSectionNames(fileName string) ([]string, error) {
	doc, e := LoadDocument(fileName)
	if e != nil {
		return nil, e
	}

	// Extract and sort valid sections names
	length := len(doc)
	names := make([]string, 0, length)
	for key := range doc {
		keyStr := key.(string)
		if keyStr != BaseKey {
			names = append(names, keyStr)
		}
	}

	return names, nil
}

func LoadDocument(fileName string) (Document, error) {
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
