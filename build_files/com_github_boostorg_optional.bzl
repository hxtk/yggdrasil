cc_library(
    name = "optional",
    includes = ["include/"],
    hdrs = glob(["include/**/*.hpp"]),
    visibility = ["//visibility:public"],
)
