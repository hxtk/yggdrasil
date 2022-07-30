load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "pgxsqli",
    srcs = ["analyzer.go"],
    deps = [
	"@com_github_praetorian_inc_gokart//util:go_default_library",
        "@org_golang_x_tools//go/analysis:go_default_library",
        "@org_golang_x_tools//go/analysis/passes/buildssa:go_default_library",
    ],
    importpath = "github.com/hxtk/yggdrasil/analyzers/pgxsqli",
    visibility = ["//visibility:public"],
)
