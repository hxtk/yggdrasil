load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/toolproxy/client/cmd/cancel",
    deps = [
        "//common/config:go_default_library",
        "//common/config/tlsconfig:go_default_library",
	"//toolproxy/client/pkg/rpc:go_default_library",
	"@com_github_sirupsen_logrus//:go_default_library",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
    ],
    visibility = [
        "//toolproxy/client/cmd:__pkg__",
    ]
)
