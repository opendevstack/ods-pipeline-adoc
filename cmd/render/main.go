// Example invocation:
//
//		go run github.com/opendevstack/ods-pipeline-adoc/cmd/render \
//			-template=sample.adoc.tmpl \
//	        -output-dir=rendered \
//			-data-source=records:sample-artifacts/*/*.json \
//			-data-source=records:sample-artifacts/*/*.yaml
//
// Parsing of data is only supported for .json and .y(a)ml files.
package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	templateGlob := flag.String("template", "", "Glob pattern from where to source templates")
	outputDir := flag.String("output-dir", "", "Output directory where to place the rendered files")
	var dataSourceFlags multiFlag
	flag.Var(&dataSourceFlags, "data-source", "Glob pattern from where to source data (may be specified multiple times)")
	var setFlags multiFlag
	flag.Var(&setFlags, "set", "Set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	flag.Parse()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if err := render(wd, *templateGlob, *outputDir, dataSourceFlags, setFlags); err != nil {
		log.Fatal(err)
	}
}
