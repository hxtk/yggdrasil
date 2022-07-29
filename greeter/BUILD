load("@io_bazel_rules_go//go:def.bzl", "go_binary")
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

go_binary(
    name = "greeter",
    srcs = ["main.go"],
    importpath = "github.com/hxtk/yggdrasil/greeter",
    visibility = ["//visibility:public"],
    deps = [
        ":greeter_go_proto",
        "//common/authz:go_default_library",
        "//common/config:go_default_library",
        "//common/server:go_default_library",
	"@org_golang_google_grpc//:go_default_library",
    ],
)

proto_library(
    name = "greeter_proto",
    srcs = ["greeter.proto"],
    deps = [
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_protobuf//:field_mask_proto",
        "@com_google_protobuf//:timestamp_proto",
	"//common/authz/v1alpha1:annotations_proto",
    ],
)

go_proto_library(
    name = "greeter_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
	#"//common/authz/protoc-gen-go-grpc-permissions:go_proto_permissions",
	#"//common/urn/protoc-gen-urn-parser:go_proto_urn_parser",
    ],
    importpath = "google.golang.org/grpc/examples/helloworld/helloworld",
    proto = ":greeter_proto",
    deps = [
        "@com_google_googleapis//google/api:annotations_go_proto",
	"@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
	#"//common/authz/v1alpha1:go_default_library",
        #"//common/authz:go_default_library",
        #"//common/urn:go_default_library",
    ],
)
