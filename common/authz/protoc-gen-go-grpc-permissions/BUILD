load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@io_bazel_rules_go//proto:compiler.bzl", "go_proto_compiler")

go_library(
    name = "lib",
    srcs = ["main.go"],
    importpath = "github.com/hxtk/yggdrasil/common/authz/protoc-gen-go-grpc-permissions",
    visibility = ["//visibility:private"],
    deps = [
        "//common/authz/protoc-gen-go-grpc-permissions/plugin",
        "@org_golang_google_protobuf//compiler/protogen",
    ],
)

go_binary(
    name = "protoc-gen-go-grpc-permissions",
    embed = [":lib"],
    visibility = ["//visibility:public"],
)

go_proto_compiler(
    name = "go_proto_permissions",
    plugin = ":protoc-gen-go-grpc-permissions",
    suffix = ".permissions.pb.go",
    valid_archive = False,
    visibility = ["//visibility:public"],
    deps = ["//common/authz"],
)
