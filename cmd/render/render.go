package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func render(workingDir, baseDir, templateGlob, outputDir string, setFlags []string) error {
	if !strings.HasSuffix(workingDir, "/") {
		workingDir = workingDir + "/"
	}
	if !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(workingDir, outputDir)
	}
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	data := map[string]any{}
	for _, v := range setFlags {
		key, value, found := strings.Cut(v, "=")
		if !found {
			return fmt.Errorf("%q is not a valid --set flag, must be of form key=value", v)
		}
		data[key] = value
	}

	matches, err := filepath.Glob(filepath.Join(workingDir, templateGlob))
	if err != nil {
		return err
	}
	for _, templateFile := range matches {
		log.Printf(
			"Rendering template %q into %q ...",
			strings.TrimPrefix(templateFile, workingDir),
			strings.TrimPrefix(outputDir, workingDir),
		)
		templateBase := filepath.Base(templateFile)
		tmpl, err := template.
			New(templateBase).
			Funcs(templateFuncs(baseDir)).
			Funcs(sprig.FuncMap()).
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
