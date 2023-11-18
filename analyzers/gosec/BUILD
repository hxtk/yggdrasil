load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "gosec",
    srcs = ["gosec.go"],
    importpath = "github.com/hxtk/yggdrasil/analyzers/gosec",
    visibility = ["//visibility:public"],
    deps = [
        "//analyzers/pkg/result",
        "@com_github_securego_gosec//:gosec",
        "@com_github_securego_gosec//rules",
        "@org_golang_x_tools//go/analysis",
        "@org_golang_x_tools//go/packages",
    ],
)
