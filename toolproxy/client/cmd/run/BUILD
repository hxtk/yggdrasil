load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "run",
    srcs = ["run.go"],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/client/cmd/run",
    visibility = [
        "//toolproxy/client/cmd:__pkg__",
    ],
    deps = [
        "//common/config/tlsconfig",
        "//toolproxy/client/pkg/rpc",
        "@com_github_sirupsen_logrus//:logrus",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
    ],
)
