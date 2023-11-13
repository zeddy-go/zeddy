// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.4
// source: error.proto

package errx

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Map struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fields map[string]*Value `protobuf:"bytes,1,rep,name=Fields,proto3" json:"Fields,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Map) Reset() {
	*x = Map{}
	if protoimpl.UnsafeEnabled {
		mi := &file_error_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Map) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Map) ProtoMessage() {}

func (x *Map) ProtoReflect() protoreflect.Message {
	mi := &file_error_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Map.ProtoReflect.Descriptor instead.
func (*Map) Descriptor() ([]byte, []int) {
	return file_error_proto_rawDescGZIP(), []int{0}
}

func (x *Map) GetFields() map[string]*Value {
	if x != nil {
		return x.Fields
	}
	return nil
}

type List struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []*Value `protobuf:"bytes,1,rep,name=List,proto3" json:"List,omitempty"`
}

func (x *List) Reset() {
	*x = List{}
	if protoimpl.UnsafeEnabled {
		mi := &file_error_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *List) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*List) ProtoMessage() {}

func (x *List) ProtoReflect() protoreflect.Message {
	mi := &file_error_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use List.ProtoReflect.Descriptor instead.
func (*List) Descriptor() ([]byte, []int) {
	return file_error_proto_rawDescGZIP(), []int{1}
}

func (x *List) GetList() []*Value {
	if x != nil {
		return x.List
	}
	return nil
}

type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Kind:
	//
	//	*Value_Int32Val
	//	*Value_Int64Val
	//	*Value_StrVal
	//	*Value_MapVal
	//	*Value_ListVal
	Kind isValue_Kind `protobuf_oneof:"Kind"`
}

func (x *Value) Reset() {
	*x = Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_error_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_error_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_error_proto_rawDescGZIP(), []int{2}
}

func (m *Value) GetKind() isValue_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (x *Value) GetInt32Val() int32 {
	if x, ok := x.GetKind().(*Value_Int32Val); ok {
		return x.Int32Val
	}
	return 0
}

func (x *Value) GetInt64Val() int64 {
	if x, ok := x.GetKind().(*Value_Int64Val); ok {
		return x.Int64Val
	}
	return 0
}

func (x *Value) GetStrVal() string {
	if x, ok := x.GetKind().(*Value_StrVal); ok {
		return x.StrVal
	}
	return ""
}

func (x *Value) GetMapVal() *Map {
	if x, ok := x.GetKind().(*Value_MapVal); ok {
		return x.MapVal
	}
	return nil
}

func (x *Value) GetListVal() *List {
	if x, ok := x.GetKind().(*Value_ListVal); ok {
		return x.ListVal
	}
	return nil
}

type isValue_Kind interface {
	isValue_Kind()
}

type Value_Int32Val struct {
	Int32Val int32 `protobuf:"varint,1,opt,name=Int32Val,proto3,oneof"`
}

type Value_Int64Val struct {
	Int64Val int64 `protobuf:"varint,2,opt,name=Int64Val,proto3,oneof"`
}

type Value_StrVal struct {
	StrVal string `protobuf:"bytes,3,opt,name=StrVal,proto3,oneof"`
}

type Value_MapVal struct {
	MapVal *Map `protobuf:"bytes,4,opt,name=MapVal,proto3,oneof"`
}

type Value_ListVal struct {
	ListVal *List `protobuf:"bytes,5,opt,name=ListVal,proto3,oneof"`
}

func (*Value_Int32Val) isValue_Kind() {}

func (*Value_Int64Val) isValue_Kind() {}

func (*Value_StrVal) isValue_Kind() {}

func (*Value_MapVal) isValue_Kind() {}

func (*Value_ListVal) isValue_Kind() {}

var File_error_proto protoreflect.FileDescriptor

var file_error_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x8c, 0x01, 0x0a, 0x03,
	0x4d, 0x61, 0x70, 0x12, 0x35, 0x0a, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x2e, 0x4d, 0x61, 0x70, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x1a, 0x4e, 0x0a, 0x0b, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x29, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x2f, 0x0a, 0x04, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x27, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x22, 0xc2, 0x01, 0x0a, 0x05,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1c, 0x0a, 0x08, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x08, 0x49, 0x6e, 0x74, 0x33, 0x32,
	0x56, 0x61, 0x6c, 0x12, 0x1c, 0x0a, 0x08, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x08, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61,
	0x6c, 0x12, 0x18, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x56, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x06, 0x53, 0x74, 0x72, 0x56, 0x61, 0x6c, 0x12, 0x2b, 0x0a, 0x06, 0x4d,
	0x61, 0x70, 0x56, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x4d, 0x61, 0x70, 0x48, 0x00,
	0x52, 0x06, 0x4d, 0x61, 0x70, 0x56, 0x61, 0x6c, 0x12, 0x2e, 0x0a, 0x07, 0x4c, 0x69, 0x73, 0x74,
	0x56, 0x61, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x48, 0x00, 0x52,
	0x07, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x42, 0x06, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64,
	0x42, 0x25, 0x5a, 0x23, 0x31, 0x31, 0x31, 0x2e, 0x32, 0x33, 0x31, 0x2e, 0x34, 0x34, 0x2e, 0x34,
	0x32, 0x2f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_error_proto_rawDescOnce sync.Once
	file_error_proto_rawDescData = file_error_proto_rawDesc
)

func file_error_proto_rawDescGZIP() []byte {
	file_error_proto_rawDescOnce.Do(func() {
		file_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_error_proto_rawDescData)
	})
	return file_error_proto_rawDescData
}

var file_error_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_error_proto_goTypes = []interface{}{
	(*Map)(nil),   // 0: common.error.Map
	(*List)(nil),  // 1: common.error.List
	(*Value)(nil), // 2: common.error.Value
	nil,           // 3: common.error.Map.FieldsEntry
}
var file_error_proto_depIdxs = []int32{
	3, // 0: common.error.Map.Fields:type_name -> common.error.Map.FieldsEntry
	2, // 1: common.error.List.List:type_name -> common.error.Value
	0, // 2: common.error.Value.MapVal:type_name -> common.error.Map
	1, // 3: common.error.Value.ListVal:type_name -> common.error.List
	2, // 4: common.error.Map.FieldsEntry.value:type_name -> common.error.Value
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_error_proto_init() }
func file_error_proto_init() {
	if File_error_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_error_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Map); i {
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
		file_error_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*List); i {
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
		file_error_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Value); i {
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
	file_error_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Value_Int32Val)(nil),
		(*Value_Int64Val)(nil),
		(*Value_StrVal)(nil),
		(*Value_MapVal)(nil),
		(*Value_ListVal)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_error_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_error_proto_goTypes,
		DependencyIndexes: file_error_proto_depIdxs,
		MessageInfos:      file_error_proto_msgTypes,
	}.Build()
	File_error_proto = out.File
	file_error_proto_rawDesc = nil
	file_error_proto_goTypes = nil
	file_error_proto_depIdxs = nil
}
