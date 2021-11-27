load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/common/authz",
    deps = [
        "//common/authz/v1alpha1:go_default_library",
        "//common/authn:go_default_library",
	"@com_authzed_api//:v1_go_proto",
	"@org_golang_google_grpc//:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
    ],
    visibility = [
        "//visibility:public",
    ],
)

cc_library(
    name = "cc_interceptor",
    srcs = glob(["*.cc"]),
    hdrs = glob(["*.h"]),
    deps = [
        "@com_github_grpc_grpc//:grpc++",
	"@ory_keto_v1alpha1//:acl_cc_grpc",
	"@com_authzed_api//:permission_service_cc_grpc",
	"@com_github_google_glog//:glog",
    ],
    visibility = ["//visibility:public"],
)
