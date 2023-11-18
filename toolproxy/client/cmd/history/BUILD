load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "history",
    srcs = ["history.go"],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/client/cmd/history",
    visibility = [
        "//toolproxy/client/cmd:__pkg__",
    ],
    deps = ["@com_github_spf13_cobra//:cobra"],
)
