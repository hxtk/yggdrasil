// Copyright (c) 2021. Peter Sanders.

#ifndef YGGDRASIL_COMMON_GRPC_LOGGING_INTERCEPTOR_H_
#define YGGDRASIL_COMMON_GRPC_LOGGING_INTERCEPTOR_H_

#include <grpcpp/grpcpp.h>
#include <grpcpp/security/auth_metadata_processor.h>

namespace yggdrasil {
namespace common {
namespace rpc {

class LoggingInterceptor: public grpc::AuthMetadataProcessor {
 public:
  virtual grpc::Status Process(
      const InputMetadata& auth_metadata,
      grpc::AuthContext* context,
      OutputMetadata* consumed_auth_metadata,
      OutputMetadata* response_metadata) override;

  inline bool IsBlocking() const override {
    return true;
  }
};

}  // namespace rpc
}  // namespace common
}  // namespace yggdrasil

#endif  // YGGDRASIL_COMMON_GRPC_LOGGING_INTERCEPTOR_H_
