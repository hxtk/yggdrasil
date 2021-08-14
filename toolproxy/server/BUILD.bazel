load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_pkg//:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_docker//container:container.bzl", "container_layer", "container_image", "container_push")

go_library(
    name = "lib",
    srcs = ["main.go"],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/server",
    visibility = ["//visibility:private"],
    deps = ["//toolproxy/server/cmd:cmd"],
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

pkg_tar(
    name = "main_config",
    srcs = ["prod.yaml"],
    extension = "tar.gz",
    mode = "0644",
    package_dir = "/etc/tool-server",
)

pkg_tar(
    name = "prod_config",
    srcs = ["prod.yaml"],
    extension = "tar.gz",
    mode = "0644",
    package_dir = "/etc/tool-server",
)

pkg_tar(
    name = "dev_config",
    srcs = ["dev.yaml"],
    extension = "tar.gz",
    mode = "0644",
    package_dir = "/etc/tool-server",
)
