load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_pkg//:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_docker//container:container.bzl", "container_layer", "container_image", "container_push")

go_library(
    name = "lib",
    srcs = ["main.go"],
    importpath = "github.com/hxtk/yggdrasil/pinger/server",
    visibility = ["//visibility:private"],
    deps = ["//pinger/server/cmd:cmd"],
)

go_binary(
    name = "pinger",
    embed = [":lib"],
    visibility = ["//visibility:public"],
)

pkg_tar(
    name = "pinger_pkg",
    srcs = [":pinger"],
    extension = "tar.gz",
    package_dir = "/usr/bin",
)

