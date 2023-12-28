load("@bazel_gazelle//:def.bzl", "gazelle")
load("@com_github_hxtk_rules_anchore//anchore:defs.bzl", "grype_updater")
load("@io_bazel_rules_go//go:def.bzl", "TOOLS_NOGO", "nogo")
load("//analyzers/gokart:gokart.bzl", _GOKART_ANALYZERS = "ANALYZERS")

exports_files(glob(["*"]))

nogo(
    name = "my_nogo",
    config = "nogo_config.json",
    visibility = ["//visibility:public"],
    deps = TOOLS_NOGO + _GOKART_ANALYZERS + [
        "//analyzers/gosec",
        "//analyzers/pgxsqli",
    ],
)

# gazelle:go_grpc_compilers @io_bazel_rules_go//proto:go_grpc,//common/authz/protoc-gen-go-grpc-permissions:go_proto_permissions,//common/urn/protoc-gen-urn-parser:go_proto_urn_parser

# gazelle:resolve proto proto google/api/client.proto @googleapis//google/api:client_proto
# gazelle:resolve proto proto google/api/annotations.proto @googleapis//google/api:annotations_proto
# gazelle:resolve proto proto google/api/field_behavior.proto @googleapis//google/api:field_behavior_proto
# gazelle:resolve proto proto google/api/resource.proto @googleapis//google/api:resource_proto
# gazelle:resolve proto proto google/type/money.proto @googleapis//google/type:money_proto
# gazelle:resolve proto go google/type/money.proto @org_golang_google_genproto//googleapis/type/money
# gazelle:resolve_regexp proto go google/api/[^/]+.proto @org_golang_google_genproto//googleapis/api/annotations

# gazelle:resolve proto proto google/rpc/status.proto @googleapis//google/rpc:status_proto
# gazelle:resolve proto go google/rpc/status.proto  @org_golang_google_genproto//googleapis/rpc/status
# gazelle:resolve proto google/longrunning/operations.proto @googleapis//google/longrunning:operations_proto
# gazelle:resolve proto go google/longrunning/operations.proto @com_google_cloud_go_longrunning//autogen/longrunningpb

# gazelle:prefix github.com/hxtk/yggdrasil
gazelle(name = "gazelle")

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)

grype_updater(
    name = "update-grype",
    output = "deps.bzl#grype_db",
    repository_name = "cve_database",
)
