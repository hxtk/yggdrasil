load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/common/authz/protoc-gen-go-grpc-permissions/plugin",
    deps = [
        ":resource_go_proto",
	"@org_golang_google_protobuf//compiler/protogen:go_default_library",
	"@org_golang_google_protobuf//proto:go_default_library",
	"@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
	"@org_golang_google_protobuf//types/descriptorpb:go_default_library",
    ],
    visibility = [
        "//common/urn/protoc-gen-urn-parser:__pkg__",
    ],
)

go_proto_library(
    name = "resource_go_proto",
    importpath = "google.golang.org/genproto/googleapis/api/annotations",
    proto = "@com_google_googleapis//google/api:resource_proto",
)
