load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "cmd",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/toolproxy/client/cmd",
    visibility = ["//visibility:public"],
    deps = [
        "//common/config:go_default_library",
        "//common/config/flags:go_default_library",
        "//common/config/tlsconfig:go_default_library",
	"//toolproxy/client/cmd/cancel:go_default_library",
	"//toolproxy/client/cmd/history:go_default_library",
	"//toolproxy/client/cmd/run:go_default_library",
	"//toolproxy/client/pkg/rpc:go_default_library",
        "@com_github_mitchellh_go_homedir//:go-homedir",
	"@com_github_sirupsen_logrus//:go_default_library",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
    ],
)
