load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "cmd",
    srcs = [
        "approve.go",
        "root.go",
    ],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/client/cmd",
    visibility = ["//visibility:public"],
    deps = [
        "//toolproxy/client/cmd/cancel",
        "//toolproxy/client/cmd/history",
        "//toolproxy/client/cmd/run",
        "@com_github_mitchellh_go_homedir//:go-homedir",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
    ],
)
