load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["parser.go"],
    importpath = "github.com/hxtk/yggdrasil/common/urn",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = glob(["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/common/urn",
    embed = [":go_default_library"],
)
