workspace(name = "yggdrasil")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

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
    sha256 = "528927e398f4e290001886894dac17c5c6a2e5548f3fb68004cfb01af901b53a",
    strip_prefix = "protobuf-3.17.3",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/refs/tags/v3.17.3.zip"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
    name = "com_google_googleapis",
    commit = "355a80d7efb196482e07c8c8ec249b0fdf7d3ff3",
    remote = "https://github.com/googleapis/googleapis.git",
    shallow_since = "1616721465 -0700",
)

load("@com_google_googleapis//:repository_rules.bzl", "switched_rules_by_language")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
    cc = True,
    go = True,
    grpc = True,
)

################################################################################
# Golang Rules

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "2b1641428dff9018f9e85c0384f03ec6c10660d935b750e3fa1492a281a53b0f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.29.0/rules_go-v0.29.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.29.0/rules_go-v0.29.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "de69a09dc70417580aabf20a28619bb3ef60d038470c7cf8442fafcf627c21cb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

go_repository(
    name = "org_golang_google_grpc",
    build_file_proto_mode = "disable",
    importpath = "google.golang.org/grpc",
    sum = "h1:Eeu7bZtDZ2DpRCsLhUlcrLnvYaMK1Gz86a+hMVvELmM=",
    version = "v1.43.0",
)

load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

go_rules_dependencies()

go_register_toolchains(version = "1.17.2")

gazelle_dependencies()

################################################################################
# Container Rules

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "59536e6ae64359b716ba9c46c39183403b01eabfbd57578e84398b4829ca499a",
    strip_prefix = "rules_docker-0.22.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.22.0/rules_docker-v0.22.0.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
)

container_pull(
    name = "distroless_base",
    digest = "sha256:02f4c952f790848aa6ffee8d241c67e0ac5364931c76a80094348af386076ad4",
    registry = "gcr.io",
    tag = "nonroot",
    repository = "distroless/base-debian11",
)

container_pull(
    name = "distroless_static",
    digest = "sha256:213a6d5205aa1421bd128b0396232a22fbb4eec4cbe510118f665398248f6d9a",
    registry = "gcr.io",
    tag = "nonroot",
    repository = "distroless/static-debian11",
)

################################################################################
# gRPC C++ Rules

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_python",
    sha256 = "cd6730ed53a002c56ce4e2f396ba3b3be262fd7cb68339f0377a45e8227fe332",
    url = "https://github.com/bazelbuild/rules_python/releases/download/0.5.0/rules_python-0.5.0.tar.gz",
)

http_archive(
    name = "com_github_grpc_grpc",
    sha256 = "b2f2620c762427bfeeef96a68c1924319f384e877bc0e084487601e4cc6e434c",
    strip_prefix = "grpc-1.42.0",
    urls = [
        "https://github.com/grpc/grpc/archive/v1.42.0.tar.gz",
    ],
)

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@com_github_grpc_grpc//bazel:grpc_extra_deps.bzl", "grpc_extra_deps")

grpc_extra_deps()

load("@rules_python//python/legacy_pip_import:pip.bzl", "pip_import")

pip_import(
    name = "grpc_python_dependencies",
    requirements = "@com_github_grpc_grpc//:requirements.bazel.txt",
)

load("@io_bazel_rules_python//python:pip.bzl", "pip_repositories")

pip_repositories()

load("@grpc_python_dependencies//:requirements.bzl", "pip_install")

pip_install()

################################################################################
# Ory Keto

load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")

new_git_repository(
    name = "ory_keto_v1alpha1",
    build_file = "//third_party:ory_keto_v1alpha1.bazel",
    commit = "3099ead2ef569e889e47c04204337639c89b1bf8",
    remote = "https://github.com/ory/keto.git",
    shallow_since = "1624432697 +0000",
    strip_prefix = "proto",
)

################################################################################
# protoc-gen-validate

http_archive(
    name = "com_envoyproxy_protoc_gen_validate",
    strip_prefix = "protoc-gen-validate-0.6.2",
    url = "https://github.com/envoyproxy/protoc-gen-validate/archive/refs/tags/v0.6.2.tar.gz",
)

################################################################################
# Authzed SpiceDB client APIs

new_git_repository(
    name = "com_authzed_api",
    build_file = "//third_party:com_authzed_api.bazel",
    commit = "29e93779606dac06b0eef40a8bc3e10c5267c552",
    remote = "https://github.com/authzed/api.git",
    shallow_since = "1637166970 -0500",
)

################################################################################
# C++ Dependencies

http_archive(
    name = "com_github_gflags_gflags",
    sha256 = "34af2f15cf7367513b352bdcd2493ab14ce43692d2dcd9dfc499492966c64dcf",
    strip_prefix = "gflags-2.2.2",
    urls = ["https://github.com/gflags/gflags/archive/v2.2.2.tar.gz"],
)

http_archive(
    name = "com_github_google_glog",
    sha256 = "21bc744fb7f2fa701ee8db339ded7dce4f975d0d55837a97be7d46e8382dea5a",
    strip_prefix = "glog-0.5.0",
    urls = ["https://github.com/google/glog/archive/v0.5.0.zip"],
)

http_archive(
    name = "com_github_jbeder_yaml_cpp",
    sha256 = "03d214d71b8bac32f684756003eb47a335fef8f8152d0894cf06e541eaf1c7f4",
    strip_prefix = "yaml-cpp-a6bbe0e50ac4074f0b9b44188c28cf00caf1a723/",
    urls = ["https://github.com/jbeder/yaml-cpp/archive/a6bbe0e50ac4074f0b9b44188c28cf00caf1a723.zip"],
)

http_archive(
    name = "com_google_googletest",
    sha256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    strip_prefix = "googletest-release-1.11.0/",
    urls = ["https://github.com/google/googletest/archive/refs/tags/release-1.11.0.zip"],
)

http_archive(
    name = "com_github_boostorg_optional",
    build_file = "@yggdrasil//third_party:com_github_boostorg_optional.bzl",
    sha256 = "39b43ba64d67da7e5a34871bfcdd9b3a1d88e943514fc2d4c7d5c9cc2b0c0355",
    strip_prefix = "/optional-optional-2021-03-10/",
    urls = ["https://github.com/boostorg/optional/archive/refs/tags/optional-2021-03-10.zip"],
)

################################################################################
# Rust Dependencies

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "rules_rust",
    sha256 = "224ebaf1156b6f2d3680e5b8c25191e71483214957dfecd25d0f29b2f283283b",
    strip_prefix = "rules_rust-a814d859845c420fd105c629134c4a4cb47ba3f8",
    urls = [
        # `main` branch as of 2021-06-15
        "https://github.com/bazelbuild/rules_rust/archive/a814d859845c420fd105c629134c4a4cb47ba3f8.tar.gz",
    ],
)

load("@rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()
