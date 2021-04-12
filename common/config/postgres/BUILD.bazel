load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/common/config/postgres",
    deps = [
        "@com_github_spf13_viper//:viper",
	"@com_github_golang_migrate_migrate_v4//:go_default_library",
	"@com_github_golang_migrate_migrate_v4//database/postgres:go_default_library",
	"@com_github_golang_migrate_migrate_v4//source/file:go_default_library",
    ],
    visibility = [
        "//visibility:public",
    ],
)
