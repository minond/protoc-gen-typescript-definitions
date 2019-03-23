package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

var (
	typmap = map[descriptor.FieldDescriptorProto_Type]string{
		// 0 is reserved for errors. Order is weird for historical reasons.
		descriptor.FieldDescriptorProto_TYPE_DOUBLE: "number",
		descriptor.FieldDescriptorProto_TYPE_FLOAT:  "number",
		// Not ZigZag encoded. Negative numbers take 10 bytes. Use TYPE_SINT64
		// if negative values are likely.
		descriptor.FieldDescriptorProto_TYPE_INT64:  "number",
		descriptor.FieldDescriptorProto_TYPE_UINT64: "number",
		// Not ZigZag encoded. Negative numbers take 10 bytes. Use TYPE_SINT32
		// if negative values are likely.
		descriptor.FieldDescriptorProto_TYPE_INT32:   "number",
		descriptor.FieldDescriptorProto_TYPE_FIXED64: "number",
		descriptor.FieldDescriptorProto_TYPE_FIXED32: "number",
		descriptor.FieldDescriptorProto_TYPE_BOOL:    "bool",
		descriptor.FieldDescriptorProto_TYPE_STRING:  "string",
		// Tag-delimited aggregate. Group type is deprecated and not supported
		// in proto3. However, Proto3 implementations should still be able to
		// parse the group wire format and treat group fields as unknown
		// fields.
		descriptor.FieldDescriptorProto_TYPE_GROUP: "void", // TODO
		// descriptor.FieldDescriptorProto_TYPE_MESSAGE is a special case
		// New in version 2.
		descriptor.FieldDescriptorProto_TYPE_BYTES:    "string",
		descriptor.FieldDescriptorProto_TYPE_UINT32:   "number",
		descriptor.FieldDescriptorProto_TYPE_ENUM:     "number",
		descriptor.FieldDescriptorProto_TYPE_SFIXED32: "number",
		descriptor.FieldDescriptorProto_TYPE_SFIXED64: "number",
		descriptor.FieldDescriptorProto_TYPE_SINT32:   "number",
		descriptor.FieldDescriptorProto_TYPE_SINT64:   "number",
	}
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	req := &plugin.CodeGeneratorRequest{}
	req.Unmarshal(data)

	var files []*plugin.CodeGeneratorResponse_File

	for _, pfile := range req.ProtoFile {
		var defs []string

		for _, message := range pfile.MessageType {
			defs = append(defs, def(*message.Name, obj(0, message.Field, message, req)))
		}

		files = append(files, &plugin.CodeGeneratorResponse_File{
			Name:    strptr(tsfilename(*pfile.Name)),
			Content: strptr(strings.TrimSpace(strings.Join(defs, "\n\n"))),
		})
	}

	res := &plugin.CodeGeneratorResponse{File: files}
	out, err := proto.Marshal(res)
	if err != nil {
		panic(err)
	}

	fmt.Print(string(out))
}

func strptr(str string) *string {
	return &str
}

// tsfilename generates a TypeScript definitions file name using a protobuf
// file name.
func tsfilename(protofilename string) string {
	var fullpath string
	subpath := strings.SplitN(protofilename, "/", 2)
	if len(subpath) > 1 {
		fullpath = subpath[1]
	} else {
		fullpath = subpath[0]
	}
	return strings.Replace(fullpath, ".proto", ".d.ts", 1)
}

// locate searches for a message type (eg .log.Log.DataEntry) in a proto
// request, looking at top-level message and nested messages as well.
func locate(typName string, req *plugin.CodeGeneratorRequest) *descriptor.DescriptorProto {
	parts := strings.Split(typName, ".")
	nested := len(parts) > 3

	for _, pfile := range req.ProtoFile {
		if *pfile.Package != parts[1] {
			continue
		}

		for _, message := range pfile.MessageType {
			if *message.Name != parts[2] {
				continue
			} else if !nested {
				return message
			}

			for _, nestedTyp := range message.NestedType {
				if *nestedTyp.Name == parts[3] {
					return nestedTyp
				}
			}
		}
	}

	return nil
}

// def will generate an exported TypeScript type declaration.
func def(name, body string) string {
	return fmt.Sprintf("export type %s = %s", name, body)
}

func tstyp(indent int, f *descriptor.FieldDescriptorProto, desc *descriptor.DescriptorProto, req *plugin.CodeGeneratorRequest) (typ string, msg string) {
	var ok bool

	if *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE && f.TypeName != nil {
		lookup := locate(*f.TypeName, req)
		if lookup != nil {
			typ = obj(indent, lookup.Field, lookup, req)
		} else {
			msg = fmt.Sprintf("FIXME unable to locate definition for %s", *f.TypeName)
			typ = "undefined"
		}
	} else if *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		msg = fmt.Sprintf("FIXME missing type name for %s", *f.Name)
		typ = "undefined"
	} else if _, ok = typmap[*f.Type]; !ok {
		msg = fmt.Sprintf("FIXME unknown protobuf type: %v", f.Type)
	} else {
		typ = typmap[*f.Type]
	}

	return typ, msg
}

// field generates a TypeScript object field type using a field descriptor.
func field(indent int, f *descriptor.FieldDescriptorProto, desc *descriptor.DescriptorProto, req *plugin.CodeGeneratorRequest) string {
	typ, msg := tstyp(indent, f, desc, req)
	if msg != "" {
		msg = "// " + msg
	}

	indentation := strings.Repeat("  ", indent)
	return strings.TrimSuffix(fmt.Sprintf("%s%s?: %s %s", indentation, *f.Name, typ, msg), " ")
}

// obj will generate a TypeScript object structure.
func obj(indent int, fields []*descriptor.FieldDescriptorProto, desc *descriptor.DescriptorProto, req *plugin.CodeGeneratorRequest) string {
	if desc.Options != nil && desc.Options.MapEntry != nil && *desc.Options.MapEntry {
		ktyp, kmsg := tstyp(indent+1, desc.Field[0], desc, req)
		vtyp, vmsg := tstyp(indent+1, desc.Field[1], desc, req)

		var comment string
		if kmsg != "" || vmsg != "" {
			comment = "// " + strings.TrimSpace(kmsg+" "+vmsg)
		}

		return strings.TrimSpace(fmt.Sprintf("Map<%s, %s> %s", ktyp, vtyp, comment))
	}

	defs := make([]string, len(fields))
	for i, f := range fields {
		defs[i] = field(indent+1, f, desc, req)
	}

	indentation := strings.Repeat("  ", indent)
	return fmt.Sprintf("{\n%s\n%s}", strings.Join(defs, "\n"), indentation)
}
