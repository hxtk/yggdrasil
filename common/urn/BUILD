load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "urn",
    srcs = ["parser.go"],
    importpath = "github.com/hxtk/yggdrasil/common/urn",
    visibility = ["//visibility:public"],
)

go_test(
    name = "urn_test",
    srcs = ["parser_test.go"],
    embed = [":urn"],
    importpath = "github.com/hxtk/yggdrasil/common/urn",
)
