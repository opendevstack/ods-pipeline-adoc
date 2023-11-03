package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/ods-pipeline-adoc/internal/testhelper"
)

func TestBuildMapPath(t *testing.T) {
	m := make(map[string]interface{})

	// build map path for first file
	p := "a/b/c/d.txt"
	p = strings.TrimSuffix(p, filepath.Ext(p))
	p = filepath.ToSlash(p)
	buildMapPath(m, p)
	if _, ok := m["a"]; !ok {
		t.Fatal("expect key a")
	}
	a := m["a"].(map[string]interface{})
	if _, ok := a["b"]; !ok {
		t.Fatal("expect key b")
	}
	b := a["b"].(map[string]interface{})
	if _, ok := b["c"]; !ok {
		t.Fatal("expect key c")
	}
	c := b["c"].(map[string]interface{})
	if _, ok := c["d"]; !ok {
		t.Fatal("expect key d")
	}

	// build map path for second file which overlaps first path
	p = "a/x.txt"
	p = strings.TrimSuffix(p, filepath.Ext(p))
	p = filepath.ToSlash(p)
	buildMapPath(m, p)
	if _, ok := a["x"]; !ok {
		t.Fatal("expect key x")
	}
	if _, ok := a["b"]; !ok {
		t.Fatal("expect key b")
	}
}

func TestSafeMapKey(t *testing.T) {
	want := "ods_artifacts_org_some_example"
	got := safeMapKey(".ods/artifacts/org.some-example")
	if want != got {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func TestRender(t *testing.T) {
	tempDir := testhelper.MkdirTempDir(t)
	defer testhelper.RmTempDir(tempDir)

	if err := render(
		"../../test/testdata/fixtures",
		"sample.adoc.tmpl",
		tempDir,
		[]string{
			".ods/artifacts/*/*.json",
			".ods/artifacts/*/*.yaml",
			".ods/artifacts/*/*.txt",
			"*.yaml",
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
		"error.adoc.tmpl",
		tempDir,
		[]string{
			".ods/artifacts/*/*.json",
			".ods/artifacts/*/*.yaml",
			".ods/artifacts/*/*.txt",
			"*.yaml",
		},
	)
	if err == nil {
		t.Error("Fixture template error.adoc.tmpl includes non-existent reference")
	} else if !strings.Contains(err.Error(), ".ods.artifacts.org_opendevstack_pipeline_go_foo.result.foo") {
		t.Errorf("Error must list valid references, got:\n%s", err)
	}
}
