load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "filelock",
    srcs = [
        "filelock.go",
        "filelock_fcntl.go",
        "filelock_other.go",
        "filelock_unix.go",
        "filelock_windows.go",
    ],
    importpath = "github.com/hxtk/yggdrasil/third_party/filelock",
    visibility = ["//visibility:public"],
)

go_test(
    name = "filelock_test",
    srcs = ["filelock_test.go"],
    deps = select({
        "@io_bazel_rules_go//go/platform:aix": [
            ":filelock",
        ],
        "@io_bazel_rules_go//go/platform:android": [
            ":filelock",
        ],
        "@io_bazel_rules_go//go/platform:darwin": [
            ":filelock",
        ],
        "@io_bazel_rules_go//go/platform:dragonfly": [
            ":filelock",
        ],
        "@io_bazel_rules_go//go/platform:freebsd": [
            ":filelock",
        ],
        "@io_bazel_rules_go//go/platform:illumos": [
            ":filelock",
        ],
        "@io_bazel_rules_go//go/platform:ios": [
            ":filelock",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            ":filelock",
        ],
        "@io_bazel_rules_go//go/platform:netbsd": [
            ":filelock",
        ],
        "@io_bazel_rules_go//go/platform:openbsd": [
            ":filelock",
        ],
        "@io_bazel_rules_go//go/platform:solaris": [
            ":filelock",
        ],
        "@io_bazel_rules_go//go/platform:windows": [
            ":filelock",
        ],
        "//conditions:default": [],
    }),
)
