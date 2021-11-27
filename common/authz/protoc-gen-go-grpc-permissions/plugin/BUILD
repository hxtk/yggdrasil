load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/common/authz/protoc-gen-go-grpc-permissions/plugin",
    deps = [
	"//common/authz/v1alpha1:go_default_library",
	"@org_golang_google_protobuf//compiler/protogen:go_default_library",
	"@org_golang_google_protobuf//proto:go_default_library",
	"@org_golang_google_protobuf//types/descriptorpb:go_default_library",
    ],
    visibility = [
        "//common/authz/protoc-gen-go-grpc-permissions:__pkg__",
    ],
)
