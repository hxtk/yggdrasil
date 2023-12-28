load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "record",
    srcs = [
        "append.go",
        "log_reader.go",
        "log_writer.go",
        "record.go",
    ],
    importpath = "github.com/hxtk/yggdrasil/common/record",
    visibility = ["//visibility:public"],
    deps = [
        "//third_party/crc",
        "//third_party/filelock",
        "@com_github_google_uuid//:uuid",
    ],
)

go_test(
    name = "record_test",
    srcs = ["record_test.go"],
    data = glob(["testdata/**"]),
    embed = [":record"],
    deps = [
        "@com_github_google_go_cmp//cmp",
        "@com_github_google_uuid//:uuid",
    ],
)
