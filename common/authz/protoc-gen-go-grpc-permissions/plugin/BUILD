load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "plugin",
    srcs = ["plugin.go"],
    importpath = "github.com/hxtk/yggdrasil/common/authz/protoc-gen-go-grpc-permissions/plugin",
    visibility = [
        "//common/authz/protoc-gen-go-grpc-permissions:__pkg__",
    ],
    deps = [
        "//common/authz/v1alpha1",
        "@org_golang_google_protobuf//compiler/protogen",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/descriptorpb",
    ],
)
