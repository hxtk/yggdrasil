// Copyright (c) 2021. Peter Sanders.

#ifndef YGGDRASIL_COMMON_GRPC_INTERCEPTOR_CHAIN_H_
#define YGGDRASIL_COMMON_GRPC_INTERCEPTOR_CHAIN_H_

#include <memory>
#include <vector>

#include <grpcpp/grpcpp.h>
#include <grpcpp/security/auth_metadata_processor.h>

namespace yggdrasil {
namespace common {
namespace rpc {

// InterceptorChain takes a vector of pointers to `grpc::AuthMetadataProcessor`s
// which will be run sequentially on each request. This is used to chain
// together Auth interceptors such that each one may have a distinct purpose,
// e.g., logging, authentication, authorization, rate-limiting.
//
// Process shall return grpc::Status::OK if and only if every child processor
// returns grpc::Status::OK. Otherwise, the chain fails fast and immediately
// returns the offending grpc::Status.
class InterceptorChain: public grpc::AuthMetadataProcessor {
 public:
  explicit InterceptorChain(std::vector<std::shared_ptr<grpc::AuthMetadataProcessor>> links)
      : links_(links) {}

  // Process implements grpc::AuthMetadataProcessor for AuthChain.
  virtual grpc::Status Process(
      const InputMetadata& auth_metadata,
      grpc::AuthContext* context,
      OutputMetadata* consumed_auth_metadata,
      OutputMetadata* response_metadata) override;

  inline bool IsBlocking() const override {
    return true;
  }
 private:
  std::vector<std::shared_ptr<grpc::AuthMetadataProcessor>> links_ = {};
};

}  // namespace rpc
}  // namespace common
}  // namespace yggdrasil

#endif  // YGGDRASIL_COMMON_GRPC_INTERCEPTOR_CHAIN_H_
