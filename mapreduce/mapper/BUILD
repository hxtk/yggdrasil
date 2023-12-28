load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "mapper",
    srcs = [
        "log_prefix.go",
        "log_recordwriter.go",
        "worker.go",
    ],
    importpath = "github.com/hxtk/yggdrasil/mapreduce/mapper",
    visibility = ["//visibility:public"],
    deps = [
        "//common/record",
        "//mapreduce/pkg/ioutil",
        "//mapreduce/v1alpha1",
        "//third_party/crc",
        "@com_github_google_btree//:btree",
        "@com_github_google_uuid//:uuid",
        "@org_golang_google_protobuf//types/known/emptypb",
    ],
)
