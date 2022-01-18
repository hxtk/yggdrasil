load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/common/authz",
    deps = [
        "//common/authz/v1alpha1:go_default_library",
        "//common/authn:go_default_library",
	"@com_authzed_api//:v1_go_proto",
	"@org_golang_google_grpc//:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//proto",
    ],
    visibility = [
        "//visibility:public",
    ],
)

proto_library(
    name = "test_proto",
    srcs = ["test.proto"],
    deps = [
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:field_mask_proto",
    ],
)

go_proto_library(
    name = "test_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
    ],
    importpath = "github.com/hxtk/yggdrasil/common/authz",
    proto = ":test_proto",
)

go_test(
    name = "go_default_test",
    srcs = glob(["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/common/authz",
    embed = [
        ":go_default_library",
        ":test_go_proto",
    ],
    deps = [
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//proto",
    ],
)
