load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "rpc",
    srcs = [
        "armada.go",
        "types.go",
    ],
    importpath = "github.com/hxtk/yggdrasil/armada/server/pkg/rpc",
    visibility = [
        "//armada/server/cmd:__pkg__",
    ],
    deps = [
        "//armada/v1alpha1",
        "//common/authn",
        "//common/authz",
        "//common/server",
        "@com_github_authzed_authzed_go//proto/authzed/api/v1:api",
        "@com_github_google_uuid//:uuid",
        "@com_github_lib_pq//:pq",
        "@com_github_sirupsen_logrus//:logrus",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
    ],
)
