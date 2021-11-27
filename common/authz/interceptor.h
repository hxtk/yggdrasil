// Copyright (c) 2021. Peter Sanders.

#ifndef YGGDRASIL_AUTHORIZATION_INTERCPTOR_H_
#define YGGDRASIL_AUTHORIZATION_INTERCPTOR_H_

#include <memory>

#include <grpcpp/grpcpp.h>
#include <grpcpp/security/auth_metadata_processor.h>

//#include "ory/keto/acl/v1alpha1/check_service.pb.h"
//#include "ory/keto/acl/v1alpha1/check_service.grpc.pb.h"

#include "authzed/api/v1/permission_service.pb.h"
#include "authzed/api/v1/permission_service.grpc.pb.h"

class AuthorizationInterceptor: public grpc::AuthMetadataProcessor {
 public:
  /*
  explicit AuthorizationInterceptor(std::shared_ptr<grpc::Channel> channel)
      : stub_(ory::keto::acl::v1alpha1::CheckService::NewStub(channel)) {}
      */
  explicit AuthorizationInterceptor(std::shared_ptr<grpc::Channel> channel)
      : stub_(authzed::api::v1::PermissionsService::NewStub(channel)) {}
  
  virtual grpc::Status Process(
      const InputMetadata& auth_metadata,
      grpc::AuthContext* context,
      OutputMetadata* consumed_auth_metadata,
      OutputMetadata* response_metadata) override;

  inline bool IsBlocking() const override {
    return true;
  }
 private:
  std::unique_ptr<authzed::api::v1::PermissionsService::Stub> stub_ = nullptr;
};

#endif  // YGGDRASIL_AUTHORIZATION_INTERCPTOR_H_
