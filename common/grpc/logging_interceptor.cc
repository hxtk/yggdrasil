#include "logging_interceptor.h"

#include <iostream>
#include <string>

#include <grpcpp/grpcpp.h>
#include <glog/logging.h>

using grpc::ClientContext;
using grpc::Status;
using grpc::StatusCode;

namespace yggdrasil {
namespace common {
namespace rpc {

Status LoggingInterceptor::Process(
      const InputMetadata& auth_metadata,
      grpc::AuthContext* context,
      OutputMetadata* consumed_auth_metadata,
      OutputMetadata* response_metadata) {
  // Get the RPC being executed.
  auto path_kv = auth_metadata.find(":path");
  if (path_kv == auth_metadata.end()) {
    return Status(StatusCode::INTERNAL, "Internal Error");
  }
  LOG(INFO) << "Request path: " << path_kv->second << std::endl;

  return Status::OK;
}

}  // namespace rpc
}  // namespace common
}  // namespace yggdrasil
