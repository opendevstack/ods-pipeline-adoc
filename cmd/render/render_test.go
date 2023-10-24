package main

import (
	"path/filepath"
	"testing"

	"github.com/opendevstack/ods-pipeline-adoc/internal/testhelper"
)

func TestRender(t *testing.T) {
	tempDir := testhelper.MkdirTempDir(t)
	defer testhelper.RmTempDir(tempDir)

	if err := render(
		"../../test/testdata/fixtures",
		"sample.adoc.tmpl",
		tempDir,
		[]string{
			"artifacts:sample-artifacts/*/*.json",
			"artifacts:sample-artifacts/*/*.yaml",
			"artifacts:sample-artifacts/*/*.txt",
			"data:*.yaml",
		},
	); err != nil {
		t.Fatal(err)
	}
	testhelper.CompareFiles(t,
		testhelper.Golden("sample.adoc"),
		filepath.Join(tempDir, "sample.adoc.out"),
	)
}

func TestRenderFailsOnMissingKeys(t *testing.T) {
	tempDir := testhelper.MkdirTempDir(t)
	defer testhelper.RmTempDir(tempDir)
	if err := render(
		"../../test/testdata/fixtures",
		"sample.adoc.tmpl",
		tempDir,
		[]string{},
	); err == nil {
		t.Fatal("Fixture template sample.adoc.tmpl requires data to be present")
	}
}
