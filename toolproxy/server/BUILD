load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_pkg//:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_docker//container:container.bzl", "container_layer", "container_image", "container_push")
load("//tools:anchore.bzl", "grype_test", "syft_sbom")

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

container_image(
    name = "tool_server_image",
    base = "@distroless_static//image",
    tars = [":tool_server_pkg"],
)

syft_sbom(
    name = "tool_server_sbom",
    image = ":tool_server_image.tar",
)

grype_test(
    name = "tool_server_image_vulnerabilities",
    sbom = ":tool_server_sbom",
    tags = [
        "manual",
        "external",
        "large",
    ],
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
