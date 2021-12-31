load("@bazel_gazelle//:def.bzl", "gazelle")
load("@com_github_hxtk_rules_anchore//anchore:defs.bzl", "grype_updater")

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
