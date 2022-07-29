load("@bazel_gazelle//:def.bzl", "gazelle")
load("@com_github_hxtk_rules_anchore//anchore:defs.bzl", "grype_updater")
load("@io_bazel_rules_go//go:def.bzl", "nogo", "TOOLS_NOGO")

nogo(
    name = "my_nogo",
    #vet = True,
    config = "nogo_config.json",
    deps = TOOLS_NOGO + [
        "//analyzers/gosec:go_default_library",
    ],
    visibility = ["//visibility:public"],
)

# gazelle:prefix github.com/hxtk/yggdrasil
gazelle(name = "gazelle")

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)

grype_updater(
    name = "update-grype",
    output = "deps.bzl#grype_db",
    repository_name = "cve_database",
)
