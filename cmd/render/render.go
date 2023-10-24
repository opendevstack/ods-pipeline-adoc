package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// fileDataMap maps a "filename" to arbitary data.
type fileDataMap map[string]map[string]interface{}

// keyedFileDataMap maps a key to a fileDataMap
type keyedFileDataMap map[string]fileDataMap

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func render(baseDir, templateGlob, outputDir string, dataSourceGlobs []string) error {
	if !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(baseDir, outputDir)
	}
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	data, err := extractDataFromSources(baseDir, dataSourceGlobs)
	if err != nil {
		return err
	}

	matches, err := filepath.Glob(filepath.Join(baseDir, templateGlob))
	if err != nil {
		return err
	}
	for _, templateFile := range matches {
		tmpl, err := template.ParseFiles(templateFile)
		if err != nil {
			return fmt.Errorf("parse template %q: %s", templateFile, err)
		}
		templateBase := filepath.Base(templateFile)
		err = renderTemplate(outputDir, templateBase, tmpl, data)
		if err != nil {
			return fmt.Errorf("render template %q: %s", templateBase, err)
		}
	}
	return nil
}

func renderTemplate(outputDir, templateFile string, tmpl *template.Template, data keyedFileDataMap) error {
	outFile := strings.TrimSuffix(templateFile, filepath.Ext(templateFile)) + ".out"
	w, err := os.Create(filepath.Join(outputDir, outFile))
	if err != nil {
		return fmt.Errorf("create output file %q: %s", outFile, err)
	}
	return tmpl.Option("missingkey=error").Execute(w, data)
}

func safeMapKey(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "_")
}

func extractDataFromSources(baseDir string, dataSourceGlobs []string) (keyedFileDataMap, error) {
	data := make(keyedFileDataMap)
	for _, keyedGlob := range dataSourceGlobs {
		globKey, globPattern, found := strings.Cut(keyedGlob, ":")
		if !found {
			return nil, fmt.Errorf("no colon found in %q", keyedGlob)
		}
		log.Printf("Collecting data from files matching glob %q under key %q ...", globPattern, globKey)
		if _, ok := data[globKey]; !ok {
			data[globKey] = make(fileDataMap)
		}
		fdm, err := collectDataFromMatchingFiles(baseDir, globPattern)
		if err != nil {
			return nil, err
		}
		for fileKey, fileData := range fdm {
			if _, ok := data[globKey][fileKey]; !ok {
				data[globKey][fileKey] = fileData
			} else {
				for k, v := range fileData {
					data[globKey][fileKey][k] = v
				}
			}
		}
	}
	return data, nil
}

func collectDataFromMatchingFiles(baseDir, glob string) (fileDataMap, error) {
	result := make(fileDataMap)
	matches, err := filepath.Glob(filepath.Join(baseDir, glob))
	if err != nil {
		return nil, err
	}
	for _, m := range matches {
		ext := filepath.Ext(m)
		decoderFunc := selectNewDecoderFunc(ext)
		matchForKey := strings.TrimPrefix(m, baseDir)
		fileName := filepath.Base(matchForKey)
		kind := filepath.Base(strings.TrimSuffix(matchForKey, fileName))
		fKey := strings.TrimSuffix(fileName, ext)

		f, err := os.Open(m)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		dec := decoderFunc(f)
		res := make(map[string]interface{})
		err = dec.Decode(&res)
		if err != nil {
			return nil, err
		}

		mk := safeMapKey(kind)
		if mk == "" || mk == "_" {
			mk = "root"
		}
		if _, ok := result[mk]; !ok {
			result[mk] = make(map[string]interface{})
		}
		result[mk][fKey] = res
	}
	return result, err
}
