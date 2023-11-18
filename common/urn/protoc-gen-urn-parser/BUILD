load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@io_bazel_rules_go//proto:compiler.bzl", "go_proto_compiler")

go_library(
    name = "lib",
    srcs = ["main.go"],
    importpath = "github.com/hxtk/yggdrasil/common/urn/protoc-gen-urn-parser",
    visibility = ["//visibility:private"],
    deps = [
        "//common/urn/protoc-gen-urn-parser/plugin",
        "@org_golang_google_protobuf//compiler/protogen",
    ],
)

go_binary(
    name = "protoc-gen-urn-parser",
    embed = [":lib"],
    visibility = ["//visibility:public"],
)

go_proto_compiler(
    name = "go_proto_urn_parser",
    plugin = ":protoc-gen-urn-parser",
    suffix = ".urn_parser.pb.go",
    valid_archive = False,
    visibility = ["//visibility:public"],
    deps = ["//common/urn"],
)
