load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "cmd",
    srcs = ["root.go"],
    importpath = "github.com/hxtk/yggdrasil/pinger/server/cmd",
    visibility = ["//visibility:public"],
    deps = [
        "//common/config/tlsconfig:go_default_library",
        "//common/grpc/server:go_default_library",
	"@com_github_sirupsen_logrus//:go_default_library",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
	"@com_github_prometheus_client_golang//prometheus:go_default_library",
	"@com_github_prometheus_client_golang//prometheus/promauto:go_default_library",
	"@com_github_tatsushid_go_fastping//:go_default_library",
    ],
)
