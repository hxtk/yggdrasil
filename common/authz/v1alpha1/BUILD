load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "annotations_proto",
    srcs = ["annotations.proto"],
    deps = [
        "@com_google_protobuf//:descriptor_proto",
    ],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "annotations_go_proto",
    importpath = "github.com/hxtk/yggdrasil/common/authz/v1alpha1",
    proto = ":annotations_proto",
)

go_library(
    name = "go_default_library",
    embed = [":annotations_go_proto"],
    importpath = "github.com/hxtk/yggdrasil/common/authz/v1alpha1",
    visibility = ["//visibility:public"],
)
