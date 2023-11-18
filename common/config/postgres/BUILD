load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "postgres",
    srcs = ["viper.go"],
    importpath = "github.com/hxtk/yggdrasil/common/config/postgres",
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "@com_github_golang_migrate_migrate_v4//:migrate",
        "@com_github_golang_migrate_migrate_v4//database/postgres",
        "@com_github_golang_migrate_migrate_v4//source/file",
        "@com_github_spf13_viper//:viper",
    ],
)
