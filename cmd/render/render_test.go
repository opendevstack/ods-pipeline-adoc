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
		"../../test/testdata/fixtures",
		"sample.adoc.tmpl",
		tempDir,
		[]string{
			"keyfoo=valbar",
			"keybar=valbaz",
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
	err := render(
		"../../test/testdata/fixtures",
		"../../test/testdata/fixtures",
		"error.adoc.tmpl",
		tempDir,
		[]string{},
	)
	if err == nil {
		t.Error("Fixture template error.adoc.tmpl includes non-existent reference, rendering should fail")
	}
}
