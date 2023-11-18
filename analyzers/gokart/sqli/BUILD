load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "sqli",
    srcs = ["analyzer.go"],
    importpath = "github.com/hxtk/yggdrasil/analyzers/gokart/sqli",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_praetorian_inc_gokart//analyzers",
        "@com_github_praetorian_inc_gokart//util",
    ],
)
