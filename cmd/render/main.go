// Example invocation:
//
//		go run github.com/opendevstack/ods-pipeline-adoc/cmd/render \
//			-template=sample.adoc.tmpl \
//	        -output-dir=rendered
package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	templateGlob := flag.String("template", "", "Glob pattern from where to source templates")
	baseDir := flag.String("base-dir", ".", "Base directory from which to interpret filepaths passed to helper functions")
	outputDir := flag.String("output-dir", "", "Output directory where to place the rendered files")
	var setFlags multiFlag
	flag.Var(&setFlags, "set", "Set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	flag.Parse()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if err := render(wd, *baseDir, *templateGlob, *outputDir, setFlags); err != nil {
		log.Fatal(err)
	}
}
