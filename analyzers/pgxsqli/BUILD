load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "pgxsqli",
    srcs = ["analyzer.go"],
    importpath = "github.com/hxtk/yggdrasil/analyzers/pgxsqli",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_praetorian_inc_gokart//util",
        "@org_golang_x_tools//go/analysis",
        "@org_golang_x_tools//go/analysis/passes/buildssa",
    ],
)
