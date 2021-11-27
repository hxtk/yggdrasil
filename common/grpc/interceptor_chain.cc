// Copyright (c) Peter Sanders. 2021.

#include "interceptor_chain.h"

#include <memory>
#include <string>
#include <vector>

#include <grpcpp/grpcpp.h>
#include <grpcpp/security/auth_metadata_processor.h>

using ::grpc::AuthContext;
using ::grpc::Status;

namespace yggdrasil {
namespace common {
namespace rpc {

Status InterceptorChain::Process(
      const InputMetadata& auth_metadata,
      AuthContext* context,
      OutputMetadata* consumed_auth_metadata,
      OutputMetadata* response_metadata) {
  for (auto it = links_.cbegin(); it != links_.cend(); ++it) {
    Status status = (*it)->Process(auth_metadata,
                                context,
                                consumed_auth_metadata,
                                response_metadata);
    if (!status.ok()) {
      return status;
    }
  }

  return Status::OK;
}

}  // namespace rpc
}  // namespace common
}  // namespace yggdrasil
