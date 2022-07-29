load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "ssrf",
    srcs = ["analyzer.go"],
    deps = [
	"@com_github_praetorian_inc_gokart//analyzers:go_default_library",
	"@com_github_praetorian_inc_gokart//util:go_default_library",
    ],
    importpath = "github.com/hxtk/yggdrasil/analyzers/gokart/ssrf",
    visibility = ["//visibility:public"],
)
