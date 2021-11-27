load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/common/server",
    deps = [
        "//common/authn:go_default_library",
        "//common/authz:go_default_library",
        "//common/config/tlsconfig:go_default_library",
	"@org_golang_google_grpc//:go_default_library",
	"@org_golang_google_grpc//reflection:go_default_library",
	"@com_github_sirupsen_logrus//:go_default_library",
	"@com_github_soheilhy_cmux//:go_default_library",
	"@com_github_grpc_ecosystem_go_grpc_middleware//:go_default_library",
	"@com_github_grpc_ecosystem_go_grpc_middleware//auth:go_default_library",
	"@com_github_grpc_ecosystem_go_grpc_middleware//logging/logrus:go_default_library",
	"@com_github_grpc_ecosystem_go_grpc_middleware//tags:go_default_library",
	"@com_github_grpc_ecosystem_go_grpc_middleware//validator:go_default_library",
	"@com_github_grpc_ecosystem_go_grpc_prometheus//:go_default_library",
	"@com_github_grpc_ecosystem_grpc_gateway//runtime:go_default_library",
	"@com_github_prometheus_client_golang//prometheus/promhttp:go_default_library",
        "@com_github_spf13_viper//:viper",
    ],
    visibility = [
        "//visibility:public"
    ],
)
