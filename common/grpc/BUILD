cc_library(
    name = "cc_interceptor",
    srcs = glob(["interceptor_chain.cc"]),
    hdrs = glob(["interceptor_chain.h"]),
    deps = [
        "@com_github_grpc_grpc//:grpc++",
    ],
    visibility = ["//visibility:public"],
)

cc_library(
    name = "logging_interceptor",
    srcs = ["logging_interceptor.cc"],
    hdrs = ["logging_interceptor.h"],
    deps = [
        "@com_github_grpc_grpc//:grpc++",
	"@com_github_google_glog//:glog",
    ],
    visibility = ["//visibility:public"],
)
