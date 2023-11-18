load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "rpc",
    srcs = ["client.go"],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/client/pkg/rpc",
    visibility = [
        "//toolproxy/client/cmd:__subpackages__",
    ],
    deps = [
        "//toolproxy/v1:toolproxy",
        "@com_github_alessio_shellescape//:shellescape",
        "@com_github_grpc_ecosystem_go_grpc_middleware//retry",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials",
    ],
)
