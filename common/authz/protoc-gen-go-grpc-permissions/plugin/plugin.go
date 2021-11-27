package plugin

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/hxtk/yggdrasil/common/authz/v1alpha1"
)

var (
	fmtPackage      = protogen.GoImportPath("fmt")
	v1alpha1Package = protogen.GoImportPath("github.com/hxtk/yggdrasil/common/authz/v1alpha1")
	authzPackage    = protogen.GoImportPath("github.com/hxtk/yggdrasil/common/authz")
	protoPackage    = protogen.GoImportPath("google.golang.org/protobuf/proto")
)

func GenerateFile(gen *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + ".permissions.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	genGeneratedHeader(gen, g)
	g.P("package ", file.GoPackageName)
	g.P()

	generateUnifiedPermissions(g, file.Proto)
	generateServicePermissions(g, file.Proto)
	generatePermCheckMethods(g)
	generateRegisterMethods(g, file.Proto)

	g.Import(fmtPackage)
	g.Import(v1alpha1Package)
}

func generatePermCheckMethods(g *protogen.GeneratedFile) {
	g.P(`func Permissions(method string) (*`, v1alpha1Package.Ident("PermissionsRule"), `, error) {`)
	g.P(`if res, ok := resourcePermissions_ALL[method]; ok {`)
	g.P(`return res, nil`)
	g.P(`}`)
	g.P(`return nil, `, fmtPackage.Ident("Errorf"), `("no permissions found for method %s", method)`)
	g.P(`}`)
	g.P()
}

func generateRegisterMethods(g *protogen.GeneratedFile, file *descriptorpb.FileDescriptorProto) {
	for _, service := range file.GetService() {
		genRegisterComment(g, service)
		g.P(`func Register`, service.GetName(), `Permissions(reg `, authzPackage.Ident(`Registrar`), `) {`)
		g.P(`for k, v := range resourcePermissions_`, service.GetName(), `{`)
		g.P(`reg.RegisterPermission(k, v)`)
		g.P(`}`)
		g.P(`}`)
		g.P()
	}
}

func genRegisterComment(g *protogen.GeneratedFile, service *descriptorpb.ServiceDescriptorProto) {
	g.P(`// Register`, service.GetName(), `Permissions registers the static permissions of `, service.GetName(), ` to reg.`)
}

func generateServicePermissions(g *protogen.GeneratedFile, file *descriptorpb.FileDescriptorProto) {

	for _, service := range file.GetService() {
		g.P(`var resourcePermissions_`, service.GetName(), ` = map[string]*`, v1alpha1Package.Ident("PermissionsRule{"))
		servicePerm, _ := extractServiceAPIOptions(service)
		for _, method := range service.GetMethod() {
			perm, err := extractMethodAPIOptions(method)
			if err != nil && servicePerm == nil {
				continue
			}

			fullMethodName := fmt.Sprintf(
				"/%s.%s/%s",
				file.GetPackage(),
				service.GetName(),
				method.GetName(),
			)
			g.P(`"`, fullMethodName, `": {`)
			if perm != nil {
				g.P(`ResourceType: "`, perm.GetResourceType(), `",`)
				g.P(`Permission:   "`, perm.GetPermission(), `",`)
			} else {
				g.P(`ResourceType: "`, servicePerm.GetResourceType(), `",`)
				g.P(`Permission:   "`, servicePerm.GetPermission(), `",`)
			}
			g.P(`},`)
		}
		g.P(`}`)
		g.P()
	}

}

func generateUnifiedPermissions(g *protogen.GeneratedFile, file *descriptorpb.FileDescriptorProto) {
	g.P(`var resourcePermissions_ALL = map[string]*`, v1alpha1Package.Ident("PermissionsRule{"))

	for _, service := range file.GetService() {
		servicePerm, _ := extractServiceAPIOptions(service)
		for _, method := range service.GetMethod() {
			perm, err := extractMethodAPIOptions(method)
			if err != nil && servicePerm == nil {
				continue
			}

			fullMethodName := fmt.Sprintf(
				"/%s.%s/%s",
				file.GetPackage(),
				service.GetName(),
				method.GetName(),
			)
			g.P(`"`, fullMethodName, `": {`)
			if perm != nil {
				g.P(`ResourceType: "`, perm.GetResourceType(), `",`)
				g.P(`Permission:   "`, perm.GetPermission(), `",`)
			} else {
				g.P(`ResourceType: "`, servicePerm.GetResourceType(), `",`)
				g.P(`Permission:   "`, servicePerm.GetPermission(), `",`)
			}
			g.P(`},`)
		}
	}
	g.P(`}`)
	g.P()

}

func genGeneratedHeader(gen *protogen.Plugin, g *protogen.GeneratedFile) {
	g.P("// Code generated by protoc-gen-gogrpc-permissions. DO NOT EDIT.")
	g.P()
}

func extractServiceAPIOptions(serv *descriptorpb.ServiceDescriptorProto) (*v1alpha1.PermissionsRule, error) {
	if serv.Options == nil {
		return nil, fmt.Errorf("no options found")
	}

	if !proto.HasExtension(serv.Options, v1alpha1.E_DefaultPermissions) {
		return nil, fmt.Errorf("no permisisons extension found")
	}

	ext := proto.GetExtension(serv.Options, v1alpha1.E_DefaultPermissions)
	opts, ok := ext.(*v1alpha1.PermissionsRule)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want a PermissionsRule", ext)
	}
	return opts, nil
}

func extractMethodAPIOptions(meth *descriptorpb.MethodDescriptorProto) (*v1alpha1.PermissionsRule, error) {
	if meth.Options == nil {
		return nil, fmt.Errorf("no options found")
	}

	if !proto.HasExtension(meth.Options, v1alpha1.E_Permissions) {
		return nil, fmt.Errorf("no permisisons extension found")
	}

	ext := proto.GetExtension(meth.Options, v1alpha1.E_Permissions)
	opts, ok := ext.(*v1alpha1.PermissionsRule)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want a PermissionsRule", ext)
	}
	return opts, nil
}
