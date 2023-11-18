load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_pkg//:pkg.bzl", "pkg_tar")
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_tarball")
load("@com_github_hxtk_rules_anchore//anchore:anchore.bzl", "grype_test")

go_library(
    name = "lib",
    srcs = ["main.go"],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/server",
    visibility = ["//visibility:private"],
    deps = ["//toolproxy/server/cmd"],
)

go_binary(
    name = "tool-server",
    embed = [":lib"],
    visibility = ["//visibility:public"],
)

pkg_tar(
    name = "tool_server_pkg",
    srcs = [":tool-server"],
    extension = "tar.gz",
    package_dir = "/usr/bin",
)

oci_image(
    name = "server_image",
    base = "@distroless_static",
    tars = [":tool_server_pkg"],
)

oci_tarball(
    name = "server_image.tar",
    image = ":server_image",
    repo_tags = ["toolproxy:latest"],
)

grype_test(
    name = "cve_scan",
    timeout = "short",
    database = "@cve_database//file",
    image = ":server_image.tar",
)
