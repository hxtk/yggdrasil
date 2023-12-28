load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "ioutil",
    srcs = [
        "reader.go",
        "writer.go",
    ],
    importpath = "github.com/hxtk/yggdrasil/mapreduce/pkg/ioutil",
    visibility = ["//visibility:public"],
    deps = ["@org_golang_google_protobuf//proto"],
)
