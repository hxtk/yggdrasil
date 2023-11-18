load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "armada_proto",
    srcs = ["armada.proto"],
    visibility = ["//visibility:public"],
    deps = [
        "//common/authz/v1alpha1:annotations_proto",
        "@com_google_protobuf//:duration_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:field_mask_proto",
        "@com_google_protobuf//:timestamp_proto",
        "@googleapis//google/api:annotations_proto",
        "@googleapis//google/api:client_proto",
        "@googleapis//google/api:field_behavior_proto",
        "@googleapis//google/api:resource_proto",
        "@googleapis//google/type:money_proto",
    ],
)

go_proto_library(
    name = "armada_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
        "//common/authz/protoc-gen-go-grpc-permissions:go_proto_permissions",
        "//common/urn/protoc-gen-urn-parser:go_proto_urn_parser",
    ],
    importpath = "github.com/hxtk/yggdrasil/armada/v1alpha1",
    proto = ":armada_proto",
    visibility = ["//visibility:public"],
    deps = [
        "//common/authz/v1alpha1",
        "@org_golang_google_genproto//googleapis/api/annotations",
        "@org_golang_google_genproto//googleapis/type/money",
    ],
)

go_library(
    name = "v1alpha1",
    srcs = ["doc.go"],
    embed = [
        ":armada_go_proto",
    ],
    importpath = "github.com/hxtk/yggdrasil/armada/v1alpha1",
    visibility = [
        "//armada/client/pkg/rpc:__pkg__",
        "//armada/server/pkg/rpc:__pkg__",
    ],
)
