load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/toolproxy/client/pkg/rpc",
    deps = [
	"//toolproxy/v1:go_default_library",
	"@org_golang_google_grpc//:go_default_library",
	"@org_golang_google_grpc//credentials:go_default_library",
	"@com_github_alessio_shellescape//:go_default_library",
	"@com_github_grpc_ecosystem_go_grpc_middleware//:go_default_library",
	"@com_github_grpc_ecosystem_go_grpc_middleware//retry:go_default_library",
	"@com_github_sirupsen_logrus//:go_default_library",
    ],
    visibility = [
        "//toolproxy/client/cmd:__subpackages__",
    ],
)
