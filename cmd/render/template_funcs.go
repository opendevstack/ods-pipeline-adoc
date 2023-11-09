package main

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"sigs.k8s.io/yaml"
)

func templateFuncs(baseDir string) template.FuncMap {
	return template.FuncMap{
		"data":          makeDataFunc(baseDir),
		"content":       makeContentFunc(baseDir),
		"contents":      makeContentsFunc(baseDir),
		"files":         makeFilesFunc(baseDir),
		"directories":   makeDirectoriesFunc(baseDir),
		"exists":        makeExistsFunc(baseDir),
		"fromMultiYAML": fromMultiYAML,
		"toYAML":        toYAML,
		"toSentence":    toSentence,
	}
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

// makeDataFunc returns a function that parses the contents of the given filename
// into a map. JSON, YAML and XML are supported.
func makeDataFunc(baseDir string) func(filename string) (map[string]any, error) {
	return func(filename string) (map[string]any, error) {
		decoderFunc := selectNewDecoderFunc(filepath.Ext(filename))

		f, err := os.Open(filepath.Join(baseDir, filename))
		if err != nil {
			return nil, err
		}
		defer f.Close()

		data := map[string]any{}
		dec := decoderFunc(f)
		err = dec.Decode(&data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

// makeContentFunc returns a function that returns the contents of the given filename.
func makeContentFunc(baseDir string) func(filename string) (string, error) {
	return func(filename string) (string, error) {
		b, err := os.ReadFile(filepath.Join(baseDir, filename))
		return strings.TrimSpace(string(b)), err
	}
}

// makeContentsFunc returns a function that returns a map with the contents of the given file glob.
func makeContentsFunc(baseDir string) func(glob string) (map[string]any, error) {
	return func(glob string) (map[string]any, error) {
		matches, err := filepath.Glob(filepath.Join(baseDir, glob))
		if err != nil {
			return nil, err
		}
		data := make(map[string]any)
		for _, m := range matches {
			f, err := os.Open(m)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			fInfo, err := f.Stat()
			if err != nil {
				return nil, err
			}
			if fInfo.IsDir() {
				continue
			}
			b, err := os.ReadFile(m)
			if err != nil {
				return nil, err
			}
			data[filepath.Base(m)] = strings.TrimSpace(string(b))
		}
		return data, nil
	}
}

// makeFilesFunc returns a function that returns the files found at given path.
func makeFilesFunc(baseDir string) func(filename string) ([]string, error) {
	return func(path string) ([]string, error) {
		entries, err := os.ReadDir(filepath.Join(baseDir, path))
		if err != nil {
			return nil, err
		}
		var dirs []string
		for _, e := range entries {
			if !e.IsDir() {
				dirs = append(dirs, e.Name())
			}
		}
		return dirs, nil
	}
}

// makeDirectoriesFunc returns a function that returns the directories found at given path.
func makeDirectoriesFunc(baseDir string) func(filename string) ([]string, error) {
	return func(path string) ([]string, error) {
		entries, err := os.ReadDir(filepath.Join(baseDir, path))
		if err != nil {
			return nil, err
		}
		var dirs []string
		for _, e := range entries {
			if e.IsDir() {
				dirs = append(dirs, e.Name())
			}
		}
		return dirs, nil
	}
}

// makeExistsFunc returns a function that checks whether the given path exists.
func makeExistsFunc(baseDir string) func(filename string) bool {
	return func(path string) bool {
		_, err := os.Open(filepath.Join(baseDir, path))
		return err == nil
	}
}
