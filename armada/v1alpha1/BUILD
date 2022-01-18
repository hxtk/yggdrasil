load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@com_googleapis_gapic_generator_go//rules_go_gapic:go_gapic.bzl", "go_gapic_library")

proto_library(
    name = "armada_proto",
    srcs = ["armada.proto"],
    deps = [
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:resource_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
        "@com_google_googleapis//google/api:client_proto",
        "@com_google_googleapis//google/type:money_proto",
        "@com_google_protobuf//:duration_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:field_mask_proto",
        "@com_google_protobuf//:timestamp_proto",
	"//common/authz/v1alpha1:annotations_proto",
    ],
)

go_proto_library(
    name = "armada_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
        "@com_github_grpc_ecosystem_grpc_gateway//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
	"//common/authz/protoc-gen-go-grpc-permissions:go_proto_permissions",
	"//common/urn/protoc-gen-urn-parser:go_proto_urn_parser",
    ],
    importpath = "github.com/hxtk/yggdrasil/armada/v1alpha1",
    proto = ":armada_proto",
    deps = [
        "@com_google_googleapis//google/api:annotations_go_proto",
	"@com_google_googleapis//google/type:money_go_proto",
	"@io_bazel_rules_go//proto/wkt:duration_go_proto",
	"@io_bazel_rules_go//proto/wkt:empty_go_proto",
	"@io_bazel_rules_go//proto/wkt:field_mask_go_proto",
	"@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
	"//common/authz/v1alpha1:go_default_library",
        "//common/authz:go_default_library",
        "//common/urn:go_default_library",
    ],
)

go_gapic_library(
    name = "armada_go_gapic",
    srcs = [
        ":armada_proto",
    ],
    deps = [
        ":armada_go_proto",
    ],
    importpath = "github.com/hxtk/yggdrasil/armada/v1alpha1/client;armada",
)

go_library(
    name = "go_default_library",
    embed = [
        ":armada_go_proto",
    ],
    importpath = "github.com/hxtk/yggdrasil/armada/v1alpha1",
    visibility = [
        "//armada/server/pkg/rpc:__pkg__",
        "//armada/client/pkg/rpc:__pkg__",
    ],
)
