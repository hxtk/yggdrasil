load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

go_library(
    name = "authz",
    srcs = ["interceptor.go"],
    embed = [":test_go_proto"],
    importpath = "github.com/hxtk/yggdrasil/common/authz",
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "//common/authn",
        "//common/authz/v1alpha1",
        "@com_github_authzed_authzed_go//proto/authzed/api/v1:api",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoreflect",
    ],
)

proto_library(
    name = "test_proto",
    srcs = ["test.proto"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:field_mask_proto",
    ],
)

go_proto_library(
    name = "test_go_proto",
    importpath = "github.com/hxtk/yggdrasil/common/authz",
    proto = ":test_proto",
    visibility = ["//visibility:public"],
)

go_test(
    name = "authz_test",
    srcs = ["interceptor_test.go"],
    embed = [":authz"],
    importpath = "github.com/hxtk/yggdrasil/common/authz",
    deps = ["@org_golang_google_protobuf//proto"],
)
