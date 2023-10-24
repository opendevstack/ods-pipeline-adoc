package e2e

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/opendevstack/ods-pipeline-adoc/internal/testhelper"
	ott "github.com/opendevstack/ods-pipeline/pkg/odstasktest"
	ttr "github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

func TestRenderAdocTask(t *testing.T) {
	if err := ttr.RunTask(
		ttr.InNamespace(namespaceConfig.Name),
		ttr.UsingTask("ods-pipeline-adoc-render"),
		ttr.WithStringParams(map[string]string{
			"template":   "templates/*.adoc.tmpl",
			"output-dir": "rendered",
		}),
		ott.WithGitSourceWorkspace(t, "../testdata/workspaces/sample-app", namespaceConfig.Name),
		ttr.AfterRun(func(config *ttr.TaskRunConfig, run *tekton.TaskRun, logs bytes.Buffer) {
			wd := config.WorkspaceConfigs["source"].Dir

			ott.AssertFilesExist(t, wd,
				"rendered/one.adoc.pdf",
				"rendered/one.adoc.out",
				"rendered/two.adoc.pdf",
				"rendered/two.adoc.out",
			)

			testhelper.CompareFiles(t,
				testhelper.Golden("sample-app/one.adoc"),
				filepath.Join(wd, "rendered/one.adoc.out"),
			)

			testhelper.CompareFiles(t,
				testhelper.Golden("sample-app/two.adoc"),
				filepath.Join(wd, "rendered/two.adoc.out"),
			)

		}),
	); err != nil {
		t.Fatal(err)
	}
}
