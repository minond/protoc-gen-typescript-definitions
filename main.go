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

	var defs []string

	for _, pfile := range req.ProtoFile {
		for _, message := range pfile.MessageType {
			defs = append(defs, def(*message.Name, obj(message.Field)))
		}
	}

	// TODO get file name using pfile.Name
	// TODO generate multiple files
	file := &plugin.CodeGeneratorResponse_File{
		Name:    strptr("definitions.d.ts"),
		Content: strptr(strings.TrimSpace(strings.Join(defs, "\n\n"))),
	}

	res := &plugin.CodeGeneratorResponse{
		File: []*plugin.CodeGeneratorResponse_File{file},
	}
	out, err := proto.Marshal(res)
	if err != nil {
		panic(err)
	}

	fmt.Print(string(out))
}

func strptr(str string) *string {
	return &str
}

// field generates a TypeScript object field type using a field descriptor.
func field(f *descriptor.FieldDescriptorProto) string {
	var msg string
	var typ string
	var ok bool

	if *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE && f.TypeName != nil {
		typ = fmt.Sprintf("any // %s", *f.TypeName)
	} else if *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		msg = fmt.Sprintf("// FIXME missing type name for %s", *f.Name)
		typ = "undefined"
	} else if _, ok = typmap[*f.Type]; !ok {
		msg = fmt.Sprintf("// FIXME unknown protobuf type: %v", f.Type)
	} else {
		typ = typmap[*f.Type]
	}

	return strings.TrimSuffix(fmt.Sprintf("  %s: %s %s", *f.Name, typ, msg), " ")
}

// def will generate an exported TypeScript type declaration.
func def(name, body string) string {
	return fmt.Sprintf("export type %s = %s", name, body)
}

// obj will generate a TypeScript object structure.
func obj(fields []*descriptor.FieldDescriptorProto) string {
	defs := make([]string, len(fields))
	for i, f := range fields {
		defs[i] = field(f)
	}

	return fmt.Sprintf("{\n%s\n}", strings.Join(defs, "\n"))
}
