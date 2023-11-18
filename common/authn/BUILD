load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "authn",
    srcs = ["interceptor.go"],
    importpath = "github.com/hxtk/yggdrasil/common/authn",
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "@com_github_authzed_authzed_go//proto/authzed/api/v1:api",
        "@com_github_google_uuid//:uuid",
        "@com_github_grpc_ecosystem_go_grpc_middleware//auth",
        "@com_github_sirupsen_logrus//:logrus",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//credentials",
        "@org_golang_google_grpc//peer",
        "@org_golang_google_grpc//status",
    ],
)
