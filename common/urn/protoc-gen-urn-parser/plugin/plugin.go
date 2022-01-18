package plugin

import (
	"fmt"
	"regexp"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	fmtPackage     = protogen.GoImportPath("fmt")
	urnPackage     = protogen.GoImportPath("github.com/hxtk/yggdrasil/common/urn")
	rePackage      = protogen.GoImportPath("regexp")
	stringsPackage = protogen.GoImportPath("strings")
	protoPackage   = protogen.GoImportPath("google.golang.org/protobuf/proto")
)

func GenerateFile(gen *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + ".urn_parser.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("package ", file.GoPackageName)
	g.P()

	generateURNResources(g, file)

	g.Import(fmtPackage)
	g.Import(rePackage)
}

func generateURNResources(g *protogen.GeneratedFile, file *protogen.File) {
	for i, msg := range file.Proto.GetMessageType() {
		generateVars(g, msg, file.Messages[i])
		generateParser(g, msg, file.Messages[i])
		generateBuilder(g, msg, file.Messages[i])
	}
}

var patternRe = regexp.MustCompile(`{([^}]+)}`)

func generateVars(g *protogen.GeneratedFile, msg *descriptorpb.DescriptorProto, pMsg *protogen.Message) {
	opt, err := extractResourceDescriptorOptions(msg)
	if err != nil {
		return
	}

	g.P("var urnKeys_", msg.GetName(), " = map[string][]string{")
	for _, v := range opt.GetPattern() {
		g.P(`"`, v, `": {`)
		for _, key := range patternRe.FindAllStringSubmatch(v, -1) {
			g.P(`"`, key[1], `",`)
		}
		g.P(`},`)
	}
	g.P("}")

	g.P("var urnREs_", msg.GetName(), " = map[string]*regexp.Regexp{")
	for _, v := range opt.GetPattern() {
		pattern := "^" + patternRe.ReplaceAllString(v, "(.+)") + "$"
		_, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Printf(
				"Skipping pattern %q for resource %q: could not be compiled.\n",
				v,
				msg.GetName(),
			)
			continue
		}
		g.P(`"`, v, `": `, rePackage.Ident("MustCompile"), `("`, pattern, `"),`)
	}
	g.P("}")
}

func generateBuilder(g *protogen.GeneratedFile, msg *descriptorpb.DescriptorProto, pMsg *protogen.Message) {
	_, err := extractResourceDescriptorOptions(msg)
	if err != nil {
		return
	}

}

func generateParser(g *protogen.GeneratedFile, msg *descriptorpb.DescriptorProto, pMsg *protogen.Message) {
	opt, err := extractResourceDescriptorOptions(msg)
	if err != nil {
		return
	}

	g.P("func Parse", msg.GetName(), "URN(name string) (*", urnPackage.Ident("URN"), ", error) {")
	g.P("for p, v := range urnREs_", msg.GetName(), " {")
	g.P("if !v.MatchString(name) {")
	g.P("continue")
	g.P("}")
	g.P()

	g.P("values := v.FindStringSubmatch(name)")
	g.P("res := ", urnPackage.Ident("Parse"), "(name)")
	g.P("res.Values = make(map[string]string, len(urnKeys_", msg.GetName(), "))")
	g.P("for i, v := range values[1:] {")
	g.P("res.Values[urnKeys_", msg.GetName(), "[p][i]] = v")
	g.P("}")
	g.P("return res, nil")

	g.P("}")
	g.P("return nil, ", fmtPackage.Ident("Errorf"), `("no pattern matched name field")`)
	g.P("}")

	nameField := "name"
	if opt.GetNameField() != "" {
		nameField = opt.GetNameField()
	}

	name, err := getNameFieldGoName(pMsg, nameField)
	if err != nil {
		fmt.Printf("Error fetching name field %q: %v.", nameField, err)
		return
	}

	g.P("func (x *", msg.GetName(), ") ResourceURN() (*", urnPackage.Ident("URN"), ", error) {")
	g.P("name := x.Get", name, "()")
	g.P("for p, v := range urnREs_", msg.GetName(), " {")
	g.P("if !v.MatchString(name) {")
	g.P("continue")
	g.P("}")
	g.P()

	g.P("values := v.FindStringSubmatch(name)")
	g.P("res := ", urnPackage.Ident("Parse"), "(name)")
	g.P("res.Values = make(map[string]string, len(urnKeys_", msg.GetName(), "))")
	g.P("for i, v := range values[1:] {")
	g.P("res.Values[urnKeys_", msg.GetName(), "[p][i]] = v")
	g.P("}")
	g.P("return res, nil")

	g.P("}")
	g.P("return nil, ", fmtPackage.Ident("Errorf"), `("no pattern matched name field")`)
	g.P("}")
}

func getNameFieldGoName(msg *protogen.Message, nameField string) (string, error) {
	for _, v := range msg.Fields {
		if v.Desc.TextName() != nameField {
			continue
		}

		if v.Desc.Kind() != protoreflect.StringKind {
			return "", fmt.Errorf("expected string; got %s", v.Desc.Kind().String())
		}

		return v.GoName, nil
	}
	return "", fmt.Errorf("no such name field")
}

func extractResourceDescriptorOptions(msg *descriptorpb.DescriptorProto) (*annotations.ResourceDescriptor, error) {
	if msg.Options == nil {
		return nil, fmt.Errorf("no options found")
	}

	if !proto.HasExtension(msg.Options, annotations.E_Resource) {
		return nil, fmt.Errorf("no resource extension found")
	}

	ext := proto.GetExtension(msg.Options, annotations.E_Resource)
	opts, ok := ext.(*annotations.ResourceDescriptor)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want a *ResourceDescriptor", ext)
	}
	return opts, nil
}
