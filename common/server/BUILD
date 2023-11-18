load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "server",
    srcs = [
        "mux.go",
        "viper.go",
    ],
    importpath = "github.com/hxtk/yggdrasil/common/server",
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "//common/authn",
        "//common/authz",
        "//common/config/tlsconfig",
        "@com_github_grpc_ecosystem_go_grpc_middleware//:go-grpc-middleware",
        "@com_github_grpc_ecosystem_go_grpc_middleware//auth",
        "@com_github_grpc_ecosystem_go_grpc_middleware//logging/logrus",
        "@com_github_grpc_ecosystem_go_grpc_middleware//tags",
        "@com_github_grpc_ecosystem_go_grpc_middleware//validator",
        "@com_github_grpc_ecosystem_go_grpc_prometheus//:go-grpc-prometheus",
        "@com_github_prometheus_client_golang//prometheus/promhttp",
        "@com_github_sirupsen_logrus//:logrus",
        "@com_github_spf13_viper//:viper",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//reflection",
    ],
)
