workspace(name = "yggdrasil")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

################################################################################
# Golang Rules

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "16e9fca53ed6bd4ff4ad76facc9b7b651a89db1689a2877d6fd7b82aa824e366",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.34.0/rules_go-v0.34.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.34.0/rules_go-v0.34.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "501deb3d5695ab658e82f6f6f549ba681ea3ca2a5fb7911154b5aa45596183fa",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.26.0/bazel-gazelle-v0.26.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.26.0/bazel-gazelle-v0.26.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

go_repository(
    name = "org_golang_google_grpc",
    build_file_proto_mode = "disable",
    importpath = "google.golang.org/grpc",
    sum = "h1:rQOsyJ/8+ufEDJd/Gdsz7HG220Mh9HAhFHRGnIjda0w=",
    version = "v1.48.0",
)

load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

go_rules_dependencies()

go_register_toolchains(
    nogo = "@yggdrasil//:my_nogo",
    version = "1.18.4",
)

gazelle_dependencies()

################################################################################
# Package Rules

http_archive(
    name = "rules_pkg",
    sha256 = "038f1caa773a7e35b3663865ffb003169c6a71dc995e39bf4815792f385d837d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_pkg/releases/download/0.4.0/rules_pkg-0.4.0.tar.gz",
        "https://github.com/bazelbuild/rules_pkg/releases/download/0.4.0/rules_pkg-0.4.0.tar.gz",
    ],
)

load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")

rules_pkg_dependencies()

################################################################################
# Protobuf Rules

http_archive(
    name = "com_google_protobuf",
    sha256 = "3d7764816081cb57752869d99b8d1c6523c054ceb19581737210a838d77403e0",
    strip_prefix = "protobuf-21.4",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/refs/tags/v21.4.zip"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
    name = "com_google_googleapis",
    commit = "c52559c85e69ea32916f8dd43e9d03bf9f695217",
    remote = "https://github.com/googleapis/googleapis.git",
    shallow_since = "1659075574 -0700",
)

load("@com_google_googleapis//:repository_rules.bzl", "switched_rules_by_language")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
    cc = True,
    go = True,
    grpc = True,
)

################################################################################
# Container Rules

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

load("@bazel_tools//tools/build_defs/repo:utils.bzl", "maybe")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
maybe(
    git_repository,
    name = "aspect_bazel_lib",
    commit = "794df714d7efbf5f2b986470428bea311f4fd772",
    remote = "https://github.com/aspect-build/bazel-lib",
    shallow_since = "1687478984 -0700",
)

maybe(
    git_repository,
    name = "rules_oci",
    remote = "https://github.com/archen/rules_oci",
    commit = "7700f626ad3439a51e9735e0fb1c2f69734e8a28",
    shallow_since = "1688397071 -0400",
)

load("@rules_oci//oci:dependencies.bzl", "rules_oci_dependencies")

rules_oci_dependencies()

load("@rules_oci//oci:repositories.bzl", "LATEST_CRANE_VERSION", "LATEST_ZOT_VERSION", "oci_register_toolchains")

oci_register_toolchains(
    name = "oci",
    crane_version = LATEST_CRANE_VERSION,
    # Uncommenting the zot toolchain will cause it to be used instead of crane for some tasks.
    # Note that it does not support docker-format images.
    # zot_version = LATEST_ZOT_VERSION,
)

load("@rules_oci//oci:pull.bzl", "oci_pull")

oci_pull(
    name = "distroless_base",
    digest = "sha256:02f4c952f790848aa6ffee8d241c67e0ac5364931c76a80094348af386076ad4",
    registry = "gcr.io",
    repository = "distroless/base-debian11",
    #tag = "nonroot",
)

oci_pull(
    name = "distroless_static",
    digest = "sha256:213a6d5205aa1421bd128b0396232a22fbb4eec4cbe510118f665398248f6d9a",
    registry = "gcr.io",
    repository = "distroless/static-debian11",
    #tag = "nonroot",
)

################################################################################
# gRPC C++ Rules

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "com_github_grpc_grpc",
    sha256 = "9b1f348b15a7637f5191e4e673194549384f2eccf01fcef7cc1515864d71b424",
    strip_prefix = "grpc-1.48.0",
    urls = [
        "https://github.com/grpc/grpc/archive/v1.48.0.tar.gz",
    ],
)

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@com_github_grpc_grpc//bazel:grpc_extra_deps.bzl", "grpc_extra_deps")

#grpc_extra_deps()

################################################################################
# protoc-gen-validate

http_archive(
    name = "com_envoyproxy_protoc_gen_validate",
    strip_prefix = "protoc-gen-validate-0.6.2",
    url = "https://github.com/envoyproxy/protoc-gen-validate/archive/refs/tags/v0.6.2.tar.gz",
)

################################################################################
# Authzed SpiceDB client APIs

load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")

new_git_repository(
    name = "com_authzed_api",
    build_file = "//third_party:com_authzed_api.bazel",
    commit = "29e93779606dac06b0eef40a8bc3e10c5267c552",
    remote = "https://github.com/authzed/api.git",
    shallow_since = "1637166970 -0500",
)

################################################################################
# Anchore Rules

http_archive(
    name = "com_github_hxtk_rules_anchore",
    strip_prefix = "rules_anchore-c97ba83b21d6498cfff3e4509cb4582bfb5924a2",
    urls = ["https://github.com/hxtk/rules_anchore/archive/c97ba83b21d6498cfff3e4509cb4582bfb5924a2.zip"],
)

load("@com_github_hxtk_rules_anchore//:deps.bzl", "anchore_deps")

anchore_deps()

load("@com_github_hxtk_rules_anchore//:extra_deps.bzl", "anchore_extra_deps")

# By default, this method configures a Go toolchain. If you have already
# configured a Go toolchain in your WORKSPACE, pass `configure_go=False`.
anchore_extra_deps(configure_go = False)

load("//:deps.bzl", "grype_db")

grype_db()
