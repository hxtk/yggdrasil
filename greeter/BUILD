load("@io_bazel_rules_go//go:def.bzl", "go_binary")
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

go_binary(
    name = "greeter",
    srcs = ["main.go"],
    importpath = "github.com/hxtk/yggdrasil/greeter",
    visibility = ["//visibility:public"],
    deps = [
        ":greeter_go_proto",
        "//common/authz:go_default_library",
        "//common/config:go_default_library",
        "//common/server:go_default_library",
	"@org_golang_google_grpc//:go_default_library",
    ],
)

proto_library(
    name = "greeter_proto",
    srcs = ["greeter.proto"],
)

go_proto_library(
    name = "greeter_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
    ],
    importpath = "google.golang.org/grpc/examples/helloworld/helloworld",
    proto = ":greeter_proto",
)
