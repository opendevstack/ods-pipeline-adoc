# ods-pipeline-adoc

[![Tests](https://github.com/opendevstack/ods-pipeline-adoc/actions/workflows/main.yaml/badge.svg)](https://github.com/opendevstack/ods-pipeline-adoc/actions/workflows/main.yaml)

Tekton task for use with [ODS Pipeline](https://github.com/opendevstack/ods-pipeline) to render asciidoc template into PDFs.

## Usage

```yaml
tasks:
- name: build
  taskRef:
    resolver: git
    params:
    - { name: url, value: https://github.com/opendevstack/ods-pipeline-adoc.git }
    - { name: revision, value: v0.1.0 }
    - { name: pathInRepo, value: tasks/render.yaml }
    workspaces:
    - { name: source, workspace: shared-workspace }
```

See the [documentation](https://github.com/opendevstack/ods-pipeline-adoc/blob/main/docs/render.adoc) for details and available parameters.

## About this repository

`docs` and `tasks` are generated directories from recipes located in `build`. See the `Makefile` target for how everything fits together.
