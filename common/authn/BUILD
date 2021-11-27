load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/common/authn",
    deps = [
        "@com_github_google_uuid//:go_default_library",
	"@com_github_grpc_ecosystem_go_grpc_middleware//auth:go_default_library",
	"@com_github_grpc_ecosystem_go_grpc_middleware//:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
	"@org_golang_google_grpc//credentials:go_default_library",
	"@org_golang_google_grpc//metadata:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
    ],
    visibility = [
        "//visibility:public",
    ],
)
