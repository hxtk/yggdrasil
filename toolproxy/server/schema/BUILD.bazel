load("@rules_pkg//:pkg.bzl", "pkg_tar")

pkg_tar(
    name = "schema",
    srcs = glob(["*.sql"]),
    extension = "tar.gz",
    package_dir = "/etc/tool-server/schema",
    visibility = [
        "//toolproxy/server:__pkg__",
    ],
)
