load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/common/config/tlsconfig",
    deps = [
        "@com_github_spf13_viper//:viper",
    ],
    visibility = [
        "//visibility:public",
    ],
)

cc_library(
    name = "tlsconfig",
    srcs = ["tlsconfig.cc"],
    hdrs = ["tlsconfig.h"],
    deps = [
        "@com_github_grpc_grpc//:grpc++",
	"@com_github_boostorg_optional//:optional",
	"@com_github_jbeder_yaml_cpp//:yaml-cpp",
	"@com_github_google_glog//:glog",
    ],
    visibility = [
        "//visibility:public",
    ],
)

