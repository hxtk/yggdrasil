load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "command_proto",
    srcs = ["command.proto"],
    deps = [
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_protobuf//:field_mask_proto",
        "@com_google_protobuf//:timestamp_proto",
	"//common/authz/v1alpha1:annotations_proto",
    ],
)

go_proto_library(
    name = "command_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
	"//common/authz/protoc-gen-go-grpc-permissions:go_proto_permissions",
    ],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/v1",
    proto = ":command_proto",
    deps = [
        "@com_google_googleapis//google/api:annotations_go_proto",
	"@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
	"//common/authz/v1alpha1:go_default_library",
        "//common/authz:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    embed = [":command_go_proto"],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/v1",
    visibility = [
        "//toolproxy/server/pkg/rpc:__pkg__",
        "//toolproxy/client/pkg/rpc:__pkg__",
    ],
)
