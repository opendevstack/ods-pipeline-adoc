package main

import (
	"html/template"
	"strings"

	"sigs.k8s.io/yaml"
)

var templateFuncs = template.FuncMap{
	"fromMultiYAML": fromMultiYAML,
	"toYAML":        toYAML,
	"toSentence":    toSentence,
	"keys":          keys,
}

// fromMultiYAML turns a string of multiple YAML documents
// (https://yaml.org/spec/1.2.2/#22-structures) into a slice of maps.
func fromMultiYAML(marshalled string) ([]map[string]interface{}, error) {
	parts := strings.Split(strings.TrimPrefix(strings.TrimSpace(marshalled), "---"), "---")
	res := []map[string]interface{}{}
	for _, p := range parts {
		v := make(map[string]interface{})
		err := yaml.Unmarshal([]byte(p), &v)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}

// toYAML turns the given object into a YAML string.
func toYAML(unmarshalled any) (string, error) {
	b, err := yaml.Marshal(unmarshalled)
	return string(b), err
}

// toSentence turns a slice into a string enumerating its items.
// The words are connected with commas, except for the last two words,
// which are connected with "and".
func toSentence(items []string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	default:
		return strings.Join(items[0:len(items)-1], ", ") + " and " + items[len(items)-1]
	}
}

// keys returns a slice of all keys in map m.
func keys(m map[string]any) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return
}
