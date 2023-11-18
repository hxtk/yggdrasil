load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

go_proto_library(
    name = "resource_go_proto",
    importpath = "google.golang.org/genproto/googleapis/api/annotations",
    proto = "@com_google_googleapis//google/api:resource_proto",
)

go_library(
    name = "plugin",
    srcs = ["plugin.go"],
    importpath = "github.com/hxtk/yggdrasil/common/urn/protoc-gen-urn-parser/plugin",
    visibility = ["//visibility:public"],
    deps = [
        ":resource_go_proto",
        "@org_golang_google_protobuf//compiler/protogen",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//types/descriptorpb",
    ],
)
