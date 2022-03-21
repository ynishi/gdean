// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: gdeanutil.proto

package gdean

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ResourceOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pattern   string `protobuf:"bytes,1,opt,name=pattern,proto3" json:"pattern,omitempty"`
	Body      string `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	Parent    string `protobuf:"bytes,3,opt,name=parent,proto3" json:"parent,omitempty"`
	Type      string `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	ChildType string `protobuf:"bytes,5,opt,name=child_type,json=childType,proto3" json:"child_type,omitempty"`
}

func (x *ResourceOption) Reset() {
	*x = ResourceOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gdeanutil_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceOption) ProtoMessage() {}

func (x *ResourceOption) ProtoReflect() protoreflect.Message {
	mi := &file_gdeanutil_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceOption.ProtoReflect.Descriptor instead.
func (*ResourceOption) Descriptor() ([]byte, []int) {
	return file_gdeanutil_proto_rawDescGZIP(), []int{0}
}

func (x *ResourceOption) GetPattern() string {
	if x != nil {
		return x.Pattern
	}
	return ""
}

func (x *ResourceOption) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *ResourceOption) GetParent() string {
	if x != nil {
		return x.Parent
	}
	return ""
}

func (x *ResourceOption) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *ResourceOption) GetChildType() string {
	if x != nil {
		return x.ChildType
	}
	return ""
}

var file_gdeanutil_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*ResourceOption)(nil),
		Field:         50000,
		Name:          "gdean.util.resource",
		Tag:           "bytes,50000,opt,name=resource",
		Filename:      "gdeanutil.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*ResourceOption)(nil),
		Field:         50001,
		Name:          "gdean.util.resource_reference",
		Tag:           "bytes,50001,opt,name=resource_reference",
		Filename:      "gdeanutil.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: ([]string)(nil),
		Field:         50000,
		Name:          "gdean.util.method_signature",
		Tag:           "bytes,50000,rep,name=method_signature",
		Filename:      "gdeanutil.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional gdean.util.ResourceOption resource = 50000;
	E_Resource = &file_gdeanutil_proto_extTypes[0]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional gdean.util.ResourceOption resource_reference = 50001;
	E_ResourceReference = &file_gdeanutil_proto_extTypes[1]
)

// Extension fields to descriptorpb.MethodOptions.
var (
	// repeated string method_signature = 50000;
	E_MethodSignature = &file_gdeanutil_proto_extTypes[2]
)

var File_gdeanutil_proto protoreflect.FileDescriptor

var file_gdeanutil_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x67, 0x64, 0x65, 0x61, 0x6e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0a, 0x67, 0x64, 0x65, 0x61, 0x6e, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x1a, 0x20, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x89, 0x01, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x63, 0x68, 0x69, 0x6c, 0x64, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x54, 0x79, 0x70, 0x65, 0x3a, 0x59, 0x0a, 0x08, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd0, 0x86, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x64, 0x65, 0x61, 0x6e, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x3a, 0x6a, 0x0a, 0x12, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x1d, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd1, 0x86, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x64, 0x65, 0x61, 0x6e, 0x2e, 0x75, 0x74, 0x69, 0x6c,
	0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x11, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x3a, 0x4b, 0x0a, 0x10, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x5f, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd0, 0x86, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0f,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x42,
	0x19, 0x5a, 0x17, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x79, 0x6e,
	0x69, 0x73, 0x68, 0x69, 0x2f, 0x67, 0x64, 0x65, 0x61, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_gdeanutil_proto_rawDescOnce sync.Once
	file_gdeanutil_proto_rawDescData = file_gdeanutil_proto_rawDesc
)

func file_gdeanutil_proto_rawDescGZIP() []byte {
	file_gdeanutil_proto_rawDescOnce.Do(func() {
		file_gdeanutil_proto_rawDescData = protoimpl.X.CompressGZIP(file_gdeanutil_proto_rawDescData)
	})
	return file_gdeanutil_proto_rawDescData
}

var file_gdeanutil_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_gdeanutil_proto_goTypes = []interface{}{
	(*ResourceOption)(nil),              // 0: gdean.util.ResourceOption
	(*descriptorpb.MessageOptions)(nil), // 1: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),   // 2: google.protobuf.FieldOptions
	(*descriptorpb.MethodOptions)(nil),  // 3: google.protobuf.MethodOptions
}
var file_gdeanutil_proto_depIdxs = []int32{
	1, // 0: gdean.util.resource:extendee -> google.protobuf.MessageOptions
	2, // 1: gdean.util.resource_reference:extendee -> google.protobuf.FieldOptions
	3, // 2: gdean.util.method_signature:extendee -> google.protobuf.MethodOptions
	0, // 3: gdean.util.resource:type_name -> gdean.util.ResourceOption
	0, // 4: gdean.util.resource_reference:type_name -> gdean.util.ResourceOption
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	3, // [3:5] is the sub-list for extension type_name
	0, // [0:3] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gdeanutil_proto_init() }
func file_gdeanutil_proto_init() {
	if File_gdeanutil_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gdeanutil_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceOption); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gdeanutil_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 3,
			NumServices:   0,
		},
		GoTypes:           file_gdeanutil_proto_goTypes,
		DependencyIndexes: file_gdeanutil_proto_depIdxs,
		MessageInfos:      file_gdeanutil_proto_msgTypes,
		ExtensionInfos:    file_gdeanutil_proto_extTypes,
	}.Build()
	File_gdeanutil_proto = out.File
	file_gdeanutil_proto_rawDesc = nil
	file_gdeanutil_proto_goTypes = nil
	file_gdeanutil_proto_depIdxs = nil
}