load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/armada/server/pkg/rpc",
    deps = [
        "//common/urn:go_default_library",
        "//common/authn:go_default_library",
        "//common/authz:go_default_library",
	"//common/server:go_default_library",
	"//armada/v1alpha1:go_default_library",
	"@com_authzed_api//:v1_go_proto",
	"@com_github_google_uuid//:go_default_library",
	"@com_github_jmoiron_sqlx//:go_default_library",
	"@com_github_lib_pq//:go_default_library",
	"@com_github_sirupsen_logrus//:go_default_library",
	"@org_golang_google_grpc//:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
	"@org_golang_google_protobuf//types/known/timestamppb:go_default_library",
    ],
    visibility = [
        "//armada/server/cmd:__pkg__",
    ],
)

go_test(
    name = "go_default_test",
    srcs = glob(["*_test.go"]),
    importpath = "github.com/hxtk/yggdrasil/armada/server/pkg/rpc",
    deps = ["@com_github_data_dog_go_sqlmock//:go_default_library"],
    embed = [":go_default_library"],
)
