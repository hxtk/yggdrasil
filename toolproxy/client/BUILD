load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "lib",
    srcs = ["main.go"],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/client",
    visibility = ["//visibility:private"],
    deps = [
        "//toolproxy/client/cmd:cmd"
    ],
)

go_binary(
    name = "toolcmd",
    embed = [":lib"],
    visibility = ["//visibility:public"],
)
