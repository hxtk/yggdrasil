load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "cmd",
    srcs = ["root.go"],
    importpath = "github.com/hxtk/yggdrasil/armada/server/cmd",
    visibility = ["//visibility:public"],
    deps = [
        "//armada/server/pkg/rpc",
        "//common/config/postgres",
        "//common/config/tlsconfig",
        "//common/server",
        "@com_github_mitchellh_go_homedir//:go-homedir",
        "@com_github_sirupsen_logrus//:logrus",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
    ],
)
