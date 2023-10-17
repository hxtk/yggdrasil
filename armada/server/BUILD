load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_pkg//:pkg.bzl", "pkg_tar")
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_tarball")
load("@com_github_hxtk_rules_anchore//anchore:anchore.bzl", "grype_test")

go_library(
    name = "lib",
    srcs = ["main.go"],
    importpath = "github.com/hxtk/yggdrasil/armada/server",
    visibility = ["//visibility:private"],
    deps = ["//armada/server/cmd:cmd"],
)

go_binary(
    name = "armada-server",
    embed = [":lib"],
    visibility = ["//visibility:public"],
)

pkg_tar(
    name = "armada_server_pkg",
    srcs = [":armada-server"],
    extension = "tar.gz",
    package_dir = "/usr/bin",
)

oci_image(
    name = "server_image",
    base = "@distroless_static",
    tars = [":armada_server_pkg"],
)

oci_tarball(
    name = "server_image.tar",
    image = ":server_image",
    repo_tags = ["armada:latest"],
)

grype_test(
    name = "cve_scan",
    image = ":server_image.tar",
    database = "@cve_database//file",
)
