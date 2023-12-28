load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "reducer",
    srcs = ["worker.go"],
    importpath = "github.com/hxtk/yggdrasil/mapreduce/reducer",
    visibility = ["//visibility:public"],
    deps = [
        "//mapreduce/pkg/ioutil",
        "//mapreduce/v1alpha1",
    ],
)
