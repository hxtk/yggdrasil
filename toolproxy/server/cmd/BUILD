load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "cmd",
    srcs = ["root.go"],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/server/cmd",
    visibility = ["//visibility:public"],
    deps = [
        "//common/config:go_default_library",
        "//common/config/postgres:go_default_library",
        "//common/config/tlsconfig:go_default_library",
        "//common/server:go_default_library",
	"//toolproxy/server/pkg/rpc:go_default_library",
        "@com_github_mitchellh_go_homedir//:go-homedir",
	"@com_github_sirupsen_logrus//:go_default_library",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
    ],
)
