load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["interfaces.go"],
    importpath = "github.com/hxtk/yggdrasil/common/config",
    visibility = [
        "//visibility:public",
    ],
    deps = [
       "@com_github_spf13_pflag//:go_default_library",
       "//common/config/tlsconfig:go_default_library",
    ],
)
