#include "interceptor.h"

#include <iostream>
#include <memory>
#include <string>

#include <grpcpp/grpcpp.h>
#include <glog/logging.h>

//#include "ory/keto/acl/v1alpha1/acl.pb.h"
//#include "ory/keto/acl/v1alpha1/check_service.pb.h"
//#include "ory/keto/acl/v1alpha1/check_service.grpc.pb.h"

#include "authzed/api/v1/permission_service.pb.h"
#include "authzed/api/v1/permission_service.grpc.pb.h"

using grpc::ClientContext;
using grpc::Status;
using grpc::StatusCode;

using authzed::api::v1::CheckPermissionRequest;
using authzed::api::v1::CheckPermissionResponse;
using authzed::api::v1::CheckPermissionResponse_Permissionship;
using authzed::api::v1::CheckPermissionResponse_Permissionship_PERMISSIONSHIP_HAS_PERMISSION;

//using ory::keto::acl::v1alpha1::CheckRequest;
//using ory::keto::acl::v1alpha1::CheckResponse;
//using ory::keto::acl::v1alpha1::Subject;

std::string to_owned_string(const grpc::string_ref& str) {
  return std::string(str.cbegin(), str.cend());
}

// Check that the relationship rpcs:/package.Service/Method#caller@user
// is allowed by the authorization service.
Status AuthorizationInterceptor::Process(
      const InputMetadata& auth_metadata,
      grpc::AuthContext* context,
      OutputMetadata* consumed_auth_metadata,
      OutputMetadata* response_metadata) {
  if (!context->IsPeerAuthenticated()) {
    LOG(INFO) << "Unauthenticated peer." << std::endl;
    return Status(StatusCode::UNAUTHENTICATED, "Unauthenticated");
  }

  /*
  CheckRequest req;
  req.set_namespace_("rpcs");
  req.set_relation("caller");

  // Get the RPC being executed.
  auto path_kv = auth_metadata.find(":path");
  if (path_kv == auth_metadata.end()) {
    return Status(StatusCode::INTERNAL, "Internal Error");
  }
  req.set_object(to_owned_string(path_kv->second));

  // Get the identity of the person running the RPC.
  req.mutable_subject()->set_id(to_owned_string(context->GetPeerIdentity().at(0)));

  // Send request to authorization service.
  ClientContext ctx;
  CheckResponse rsp;
  Status status = stub_->Check(&ctx, req, &rsp);
  if (!status.ok()) {
    // If authorization service could not be reached,
    // we allow any authenticated user.
    return Status::OK;
  }

  // Handle results of authorization request.
  if (rsp.allowed()) {
    return Status::OK;
  }
  */

  CheckPermissionRequest req;
  req.set_permission("call");

  // Get the RPC being executed.
  auto path_kv = auth_metadata.find(":path");
  if (path_kv == auth_metadata.end()) {
    return Status(StatusCode::INTERNAL, "Internal Error");
  }

  auto path = to_owned_string(path_kv->second);
  path.erase(0, 1);
  for (auto it = path.begin(); it != path.end(); ++it) {
    if (*it == '.') {
      *it = '_';
    }
  }

  auto object = req.mutable_resource();
  object->set_object_type("rpcs");
  object->set_object_id(path);

  // Get the identity of the person running the RPC.
  auto subject = req.mutable_subject()->mutable_object();
  subject->set_object_type("users");
  subject->set_object_id(to_owned_string(context->GetPeerIdentity().at(0)));

  LOG(INFO) << "Evaluating relationship: " 
      << "rpcs:" << path
      << "#call@users:" << subject->object_id() << std::endl;

  ClientContext ctx;
  CheckPermissionResponse rsp;
  Status status = stub_->CheckPermission(&ctx, req, &rsp);
  if (!status.ok()) {
    LOG(ERROR) << "Permission check failed: " << status.error_message() << std::endl;
 
    // If authorization service could not be reached,
    // we allow any authenticated user.
    return Status::OK;
  }

  if (rsp.permissionship() != CheckPermissionResponse_Permissionship_PERMISSIONSHIP_HAS_PERMISSION) {
    return Status(StatusCode::PERMISSION_DENIED, "Permission Denied");
  }

  return Status::OK;
}

