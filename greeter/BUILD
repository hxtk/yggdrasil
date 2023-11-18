load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

# gazelle:go_grpc_compilers @io_bazel_rules_go//proto:go_grpc

go_binary(
    name = "greeter",
    embed = [":greeter_lib"],
    importpath = "github.com/hxtk/yggdrasil/greeter",
    visibility = ["//visibility:public"],
)

proto_library(
    name = "greeter_proto",
    srcs = ["greeter.proto"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "helloworld_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "google.golang.org/grpc/examples/helloworld/helloworld",
    proto = ":greeter_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "greeter_lib",
    srcs = ["main.go"],
    importpath = "github.com/hxtk/yggdrasil/greeter",
    visibility = ["//visibility:private"],
    deps = [
        ":helloworld_go_proto",
        "//common/authz",
        "//common/server",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
