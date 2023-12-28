load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

# gazelle:go_grpc_compilers @io_bazel_rules_go//proto:go_grpc,//common/urn/protoc-gen-urn-parser:go_proto_urn_parser
proto_library(
    name = "mapreduce_proto",
    srcs = ["mapreduce.proto"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
        "@googleapis//google/api:annotations_proto",
        "@googleapis//google/api:field_behavior_proto",
        "@googleapis//google/longrunning:operations_proto",
    ],
)

go_proto_library(
    name = "mapreduce_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
        "//common/urn/protoc-gen-urn-parser:go_proto_urn_parser",
    ],
    importpath = "github.com/hxtk/yggdrasil/mapreduce/v1alpha1",
    proto = ":mapreduce_proto",
    visibility = ["//visibility:public"],
    deps = [
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_genproto//googleapis/api/annotations",
    ],
)

go_library(
    name = "v1alpha1",
    srcs = ["doc.go"],
    embed = [
        ":mapreduce_go_proto",
    ],
    importpath = "github.com/hxtk/yggdrasil/mapreduce/v1alpha1",
    visibility = ["//visibility:public"],
)
