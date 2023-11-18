load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "tlsconfig",
    srcs = [
        "tls_config.go",
        "viper.go",
    ],
    importpath = "github.com/hxtk/yggdrasil/common/config/tlsconfig",
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "@com_github_spf13_viper//:viper",
    ],
)
