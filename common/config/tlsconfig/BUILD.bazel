load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/common/config/tlsconfig",
    deps = [
        "@com_github_spf13_viper//:viper",
    ],
    visibility = [
        "//visibility:public",
    ],
)
