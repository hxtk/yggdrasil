load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["gosec.go"],
    deps = [
        "//analyzers/pkg/result:go_default_library",
        "@org_golang_x_tools//go/analysis:go_default_library",
        "@org_golang_x_tools//go/packages:go_default_library",
	"@com_github_securego_gosec//:go_default_library",
	"@com_github_securego_gosec//rules:go_default_library",
    ],
    importpath = "github.com/hxtk/yggdrasil/analyzers/gosec",
    visibility = ["//visibility:public"],
)
