load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "config",
    srcs = ["interfaces.go"],
    importpath = "github.com/hxtk/yggdrasil/common/config",
    visibility = [
        "//visibility:public",
    ],
    deps = ["@com_github_spf13_pflag//:pflag"],
)
