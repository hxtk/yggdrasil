load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["tls_config.go"],
    importpath = "github.com/hxtk/yggdrasil/common/config/flags",
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "@com_github_spf13_pflag//:go_default_library",
        "@com_github_thalesignite_crypto11//:go_default_library",
        "@org_golang_x_crypto//ssh/terminal:go_default_library",
        "//common/config/tlsconfig:go_default_library",
        "//common/config:go_default_library",
    ],
)
