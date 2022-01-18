load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_pkg//:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_docker//container:container.bzl", "container_layer", "container_image", "container_push")
load("@com_github_hxtk_rules_anchore//anchore:anchore.bzl", "grype_test")

go_library(
    name = "lib",
    srcs = ["main.go"],
    importpath = "github.com/hxtk/yggdrasil/armada/server",
    visibility = ["//visibility:private"],
    deps = ["//armada/server/cmd:cmd"],
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

container_image(
    name = "server_image",
    base = "@distroless_static//image",
    tars = [":tool_server_pkg"],
)

grype_test(
    name = "cve_scan",
    image = ":server_image.tar",
    database = "@cve_database//file",
)
