//cd .. && protoc -I ./proto ./proto/error.proto --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative
//cd .. && protoc -I ./proto ./proto/error.proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v5.26.1
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

type NoContent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NoContent) Reset() {
	*x = NoContent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_error_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoContent) ProtoMessage() {}

func (x *NoContent) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use NoContent.ProtoReflect.Descriptor instead.
func (*NoContent) Descriptor() ([]byte, []int) {
	return file_error_proto_rawDescGZIP(), []int{0}
}

type Map struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fields map[string]*Value `protobuf:"bytes,1,rep,name=Fields,proto3" json:"Fields,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Map) Reset() {
	*x = Map{}
	if protoimpl.UnsafeEnabled {
		mi := &file_error_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Map) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Map) ProtoMessage() {}

func (x *Map) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Map.ProtoReflect.Descriptor instead.
func (*Map) Descriptor() ([]byte, []int) {
	return file_error_proto_rawDescGZIP(), []int{1}
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
		mi := &file_error_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *List) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*List) ProtoMessage() {}

func (x *List) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use List.ProtoReflect.Descriptor instead.
func (*List) Descriptor() ([]byte, []int) {
	return file_error_proto_rawDescGZIP(), []int{2}
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
		mi := &file_error_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_error_proto_msgTypes[3]
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
	return file_error_proto_rawDescGZIP(), []int{3}
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
	0x0a, 0x0b, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0b, 0x0a,
	0x09, 0x4e, 0x6f, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x72, 0x0a, 0x03, 0x4d, 0x61,
	0x70, 0x12, 0x28, 0x0a, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x4d, 0x61, 0x70, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x1a, 0x41, 0x0a, 0x0b, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1c, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x22,
	0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x4c, 0x69,
	0x73, 0x74, 0x22, 0xa8, 0x01, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1c, 0x0a, 0x08,
	0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00,
	0x52, 0x08, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x12, 0x1c, 0x0a, 0x08, 0x49, 0x6e,
	0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x08,
	0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x12, 0x18, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x56,
	0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x53, 0x74, 0x72, 0x56,
	0x61, 0x6c, 0x12, 0x1e, 0x0a, 0x06, 0x4d, 0x61, 0x70, 0x56, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x04, 0x2e, 0x4d, 0x61, 0x70, 0x48, 0x00, 0x52, 0x06, 0x4d, 0x61, 0x70, 0x56,
	0x61, 0x6c, 0x12, 0x21, 0x0a, 0x07, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x48, 0x00, 0x52, 0x07, 0x4c, 0x69,
	0x73, 0x74, 0x56, 0x61, 0x6c, 0x42, 0x06, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x32, 0x26, 0x0a,
	0x04, 0x54, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x04, 0x54, 0x65, 0x73, 0x74, 0x12, 0x0a, 0x2e,
	0x4e, 0x6f, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x1a, 0x0a, 0x2e, 0x4e, 0x6f, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x7a, 0x65, 0x64, 0x64, 0x79, 0x2d, 0x67, 0x6f, 0x2f, 0x7a, 0x65, 0x64,
	0x64, 0x79, 0x2f, 0x65, 0x72, 0x72, 0x78, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_error_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_error_proto_goTypes = []interface{}{
	(*NoContent)(nil), // 0: NoContent
	(*Map)(nil),       // 1: Map
	(*List)(nil),      // 2: List
	(*Value)(nil),     // 3: Value
	nil,               // 4: Map.FieldsEntry
}
var file_error_proto_depIdxs = []int32{
	4, // 0: Map.Fields:type_name -> Map.FieldsEntry
	3, // 1: List.List:type_name -> Value
	1, // 2: Value.MapVal:type_name -> Map
	2, // 3: Value.ListVal:type_name -> List
	3, // 4: Map.FieldsEntry.value:type_name -> Value
	0, // 5: Test.Test:input_type -> NoContent
	0, // 6: Test.Test:output_type -> NoContent
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
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
			switch v := v.(*NoContent); i {
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
		file_error_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_error_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
	file_error_proto_msgTypes[3].OneofWrappers = []interface{}{
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
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
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
