load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "rpc",
    srcs = [
        "tool_proxy.go",
        "types.go",
    ],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/server/pkg/rpc",
    visibility = [
        "//toolproxy/server/cmd:__pkg__",
    ],
    deps = [
        "//common/authz",
        "//common/server",
        "//common/urn",
        "//toolproxy/v1:toolproxy",
        "@com_github_lib_pq//:pq",
        "@com_github_sirupsen_logrus//:logrus",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_test(
    name = "rpc_test",
    srcs = [
        "tool_proxy_create_test.go",
        "tool_proxy_delete_test.go",
        "tool_proxy_get_test.go",
    ],
    embed = [":rpc"],
    importpath = "github.com/hxtk/yggdrasil/toolproxy/server/pkg/rpc",
    deps = [
        "//toolproxy/v1:toolproxy",
        "@com_github_data_dog_go_sqlmock//:go-sqlmock",
        "@com_github_lib_pq//:pq",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)
