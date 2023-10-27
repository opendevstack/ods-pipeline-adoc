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

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func render(baseDir, templateGlob, outputDir string, dataSourceGlobs []string) error {
	if !strings.HasSuffix(baseDir, "/") {
		baseDir = baseDir + "/"
	}
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
		log.Printf(
			"Rendering template %q into %q ...",
			strings.TrimPrefix(templateFile, baseDir),
			strings.TrimPrefix(outputDir, baseDir),
		)
		templateBase := filepath.Base(templateFile)
		tmpl, err := template.
			New(templateBase).
			Funcs(templateFuncs).
			ParseFiles(templateFile)
		if err != nil {
			return fmt.Errorf("parse template %q: %s", templateFile, err)
		}
		err = renderTemplate(outputDir, templateBase, tmpl, data)
		if err != nil {
			return fmt.Errorf("render template %q: %s", templateBase, err)
		}
	}
	return nil
}

func renderTemplate(outputDir, templateFile string, tmpl *template.Template, data map[string]interface{}) error {
	outFile := strings.TrimSuffix(templateFile, filepath.Ext(templateFile)) + ".out"
	w, err := os.Create(filepath.Join(outputDir, outFile))
	if err != nil {
		return fmt.Errorf("create output file %q: %s", outFile, err)
	}
	return tmpl.Option("missingkey=error").Execute(w, data)
}

func safeMapKey(str string) string {
	return strings.Trim(nonAlphanumericRegex.ReplaceAllString(str, "_"), "_")
}

func extractDataFromSources(baseDir string, dataSourceGlobs []string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	for _, globPattern := range dataSourceGlobs {
		log.Printf("Collecting data from files matching glob %q ...", globPattern)
		err := collectDataFromMatchingFiles(baseDir, globPattern, data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// buildMapPath builds map keys in m corresponding to p.
// p is expected to be a filepath using slashes without an extension, e.g. "a/b/c/d".
func buildMapPath(m map[string]interface{}, p string) map[string]interface{} {
	elems := strings.SplitN(p, "/", 2)
	dir := safeMapKey(elems[0])
	if _, ok := m[dir]; !ok {
		m[dir] = make(map[string]interface{})
	}
	leaf := m[dir].(map[string]interface{})
	if len(elems) > 1 {
		return buildMapPath(leaf, elems[1])
	}
	return leaf
}

func collectDataFromMatchingFiles(baseDir, glob string, data map[string]interface{}) error {
	matches, err := filepath.Glob(filepath.Join(baseDir, glob))
	if err != nil {
		return err
	}
	for _, m := range matches {
		p := filepath.ToSlash(strings.TrimPrefix(m, baseDir))
		p = strings.TrimSuffix(p, filepath.Ext(p))
		fileData := buildMapPath(data, p)
		decoderFunc := selectNewDecoderFunc(filepath.Ext(m))

		f, err := os.Open(m)
		if err != nil {
			return err
		}
		defer f.Close()

		fInfo, err := f.Stat()
		if err != nil {
			return err
		}
		if fInfo.IsDir() {
			continue
		}

		dec := decoderFunc(f)
		err = dec.Decode(&fileData)
		if err != nil {
			return err
		}
	}
	return err
}
