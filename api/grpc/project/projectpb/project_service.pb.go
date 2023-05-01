// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: grpc/project/projectpb/project_service.proto

package projectpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Channel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Channel) Reset() {
	*x = Channel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_project_projectpb_project_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Channel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Channel) ProtoMessage() {}

func (x *Channel) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_project_projectpb_project_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Channel.ProtoReflect.Descriptor instead.
func (*Channel) Descriptor() ([]byte, []int) {
	return file_grpc_project_projectpb_project_service_proto_rawDescGZIP(), []int{0}
}

func (x *Channel) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Channel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Project struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int64      `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name     string     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Channels []*Channel `protobuf:"bytes,3,rep,name=channels,proto3" json:"channels,omitempty"`
}

func (x *Project) Reset() {
	*x = Project{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_project_projectpb_project_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Project) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Project) ProtoMessage() {}

func (x *Project) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_project_projectpb_project_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Project.ProtoReflect.Descriptor instead.
func (*Project) Descriptor() ([]byte, []int) {
	return file_grpc_project_projectpb_project_service_proto_rawDescGZIP(), []int{1}
}

func (x *Project) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Project) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Project) GetChannels() []*Channel {
	if x != nil {
		return x.Channels
	}
	return nil
}

type ListProjectsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Projects []*Project `protobuf:"bytes,1,rep,name=projects,proto3" json:"projects,omitempty"`
}

func (x *ListProjectsResponse) Reset() {
	*x = ListProjectsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_project_projectpb_project_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListProjectsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProjectsResponse) ProtoMessage() {}

func (x *ListProjectsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_project_projectpb_project_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListProjectsResponse.ProtoReflect.Descriptor instead.
func (*ListProjectsResponse) Descriptor() ([]byte, []int) {
	return file_grpc_project_projectpb_project_service_proto_rawDescGZIP(), []int{2}
}

func (x *ListProjectsResponse) GetProjects() []*Project {
	if x != nil {
		return x.Projects
	}
	return nil
}

var File_grpc_project_projectpb_project_service_proto protoreflect.FileDescriptor

var file_grpc_project_projectpb_project_service_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x70, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6c, 0x6f, 0x67, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2d, 0x0a, 0x07, 0x43, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x61, 0x0a, 0x07, 0x50, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x32, 0x0a, 0x08, 0x63, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x6c, 0x6f, 0x67, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x52, 0x08, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x22, 0x4a, 0x0a, 0x14,
	0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6c, 0x6f,
	0x67, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x08,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x32, 0x5f, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x0c, 0x4c, 0x69,
	0x73, 0x74, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x23, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6c, 0x6f, 0x67, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0xba, 0x01, 0x0a, 0x11, 0x63, 0x6f,
	0x6d, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6c, 0x6f, 0x67, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x42,
	0x13, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x73, 0x68, 0x69, 0x66, 0x74, 0x79, 0x31, 0x31, 0x2f, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x6c, 0x6f, 0x67, 0x2d, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x70, 0x62, 0xa2, 0x02, 0x03, 0x42, 0x47, 0x58, 0xaa, 0x02, 0x0d, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x6c, 0x6f, 0x67, 0x2e, 0x47, 0x72, 0x70, 0x63, 0xca, 0x02, 0x0d, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x6c, 0x6f, 0x67, 0x5c, 0x47, 0x72, 0x70, 0x63, 0xe2, 0x02, 0x19, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x6c, 0x6f, 0x67, 0x5c, 0x47, 0x72, 0x70, 0x63, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x6c, 0x6f, 0x67,
	0x3a, 0x3a, 0x47, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_project_projectpb_project_service_proto_rawDescOnce sync.Once
	file_grpc_project_projectpb_project_service_proto_rawDescData = file_grpc_project_projectpb_project_service_proto_rawDesc
)

func file_grpc_project_projectpb_project_service_proto_rawDescGZIP() []byte {
	file_grpc_project_projectpb_project_service_proto_rawDescOnce.Do(func() {
		file_grpc_project_projectpb_project_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_project_projectpb_project_service_proto_rawDescData)
	})
	return file_grpc_project_projectpb_project_service_proto_rawDescData
}

var file_grpc_project_projectpb_project_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_grpc_project_projectpb_project_service_proto_goTypes = []interface{}{
	(*Channel)(nil),              // 0: blocklog.grpc.Channel
	(*Project)(nil),              // 1: blocklog.grpc.Project
	(*ListProjectsResponse)(nil), // 2: blocklog.grpc.ListProjectsResponse
	(*emptypb.Empty)(nil),        // 3: google.protobuf.Empty
}
var file_grpc_project_projectpb_project_service_proto_depIdxs = []int32{
	0, // 0: blocklog.grpc.Project.channels:type_name -> blocklog.grpc.Channel
	1, // 1: blocklog.grpc.ListProjectsResponse.projects:type_name -> blocklog.grpc.Project
	3, // 2: blocklog.grpc.ProjectService.ListProjects:input_type -> google.protobuf.Empty
	2, // 3: blocklog.grpc.ProjectService.ListProjects:output_type -> blocklog.grpc.ListProjectsResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_grpc_project_projectpb_project_service_proto_init() }
func file_grpc_project_projectpb_project_service_proto_init() {
	if File_grpc_project_projectpb_project_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_project_projectpb_project_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Channel); i {
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
		file_grpc_project_projectpb_project_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Project); i {
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
		file_grpc_project_projectpb_project_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListProjectsResponse); i {
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
			RawDescriptor: file_grpc_project_projectpb_project_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_project_projectpb_project_service_proto_goTypes,
		DependencyIndexes: file_grpc_project_projectpb_project_service_proto_depIdxs,
		MessageInfos:      file_grpc_project_projectpb_project_service_proto_msgTypes,
	}.Build()
	File_grpc_project_projectpb_project_service_proto = out.File
	file_grpc_project_projectpb_project_service_proto_rawDesc = nil
	file_grpc_project_projectpb_project_service_proto_goTypes = nil
	file_grpc_project_projectpb_project_service_proto_depIdxs = nil
}
