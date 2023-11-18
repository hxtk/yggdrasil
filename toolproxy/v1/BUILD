load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "command_proto",
    srcs = ["command.proto"],
    visibility = ["//visibility:public"],
    deps = [
        "//common/authz/v1alpha1:annotations_proto",
        "@com_google_protobuf//:field_mask_proto",
        "@com_google_protobuf//:timestamp_proto",
        "@googleapis//google/api:annotations_proto",
    ],
)

go_proto_library(
    name = "command_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
        "//common/authz/protoc-gen-go-grpc-permissions:go_proto_permissions",
        "//common/urn/protoc-gen-urn-parser:go_proto_urn_parser",
    ],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/v1",
    proto = ":command_proto",
    visibility = ["//visibility:public"],
    deps = [
        "//common/authz/v1alpha1",
        "@org_golang_google_genproto//googleapis/api/annotations",
    ],
)

go_library(
    name = "toolproxy",
    embed = [":command_go_proto"],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/v1",
    visibility = [
        "//toolproxy/client/pkg/rpc:__pkg__",
        "//toolproxy/server/pkg/rpc:__pkg__",
    ],
)
