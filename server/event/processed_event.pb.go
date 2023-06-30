// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: event/processed_event.proto

package event

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EventType int32

const (
	EventType_FUNDING    EventType = 0
	EventType_STAKING    EventType = 1
	EventType_DEX        EventType = 2
	EventType_GOVERNANCE EventType = 3
)

// Enum value maps for EventType.
var (
	EventType_name = map[int32]string{
		0: "FUNDING",
		1: "STAKING",
		2: "DEX",
		3: "GOVERNANCE",
	}
	EventType_value = map[string]int32{
		"FUNDING":    0,
		"STAKING":    1,
		"DEX":        2,
		"GOVERNANCE": 3,
	}
)

func (x EventType) Enum() *EventType {
	p := new(EventType)
	*p = x
	return p
}

func (x EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_event_processed_event_proto_enumTypes[0].Descriptor()
}

func (EventType) Type() protoreflect.EnumType {
	return &file_event_processed_event_proto_enumTypes[0]
}

func (x EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventType.Descriptor instead.
func (EventType) EnumDescriptor() ([]byte, []int) {
	return file_event_processed_event_proto_rawDescGZIP(), []int{0}
}

type EventProcessedMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChainId       uint64    `protobuf:"varint,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	WalletAddress string    `protobuf:"bytes,2,opt,name=wallet_address,json=walletAddress,proto3" json:"wallet_address,omitempty"`
	EventType     EventType `protobuf:"varint,3,opt,name=event_type,json=eventType,proto3,enum=starscope.event.EventType" json:"event_type,omitempty"`
}

func (x *EventProcessedMsg) Reset() {
	*x = EventProcessedMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_processed_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventProcessedMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventProcessedMsg) ProtoMessage() {}

func (x *EventProcessedMsg) ProtoReflect() protoreflect.Message {
	mi := &file_event_processed_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventProcessedMsg.ProtoReflect.Descriptor instead.
func (*EventProcessedMsg) Descriptor() ([]byte, []int) {
	return file_event_processed_event_proto_rawDescGZIP(), []int{0}
}

func (x *EventProcessedMsg) GetChainId() uint64 {
	if x != nil {
		return x.ChainId
	}
	return 0
}

func (x *EventProcessedMsg) GetWalletAddress() string {
	if x != nil {
		return x.WalletAddress
	}
	return ""
}

func (x *EventProcessedMsg) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_FUNDING
}

var File_event_processed_event_proto protoreflect.FileDescriptor

var file_event_processed_event_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65,
	0x64, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x73,
	0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x90, 0x01, 0x0a, 0x11, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x65, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64,
	0x12, 0x25, 0x0a, 0x0e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x73, 0x74,
	0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x2a, 0x3e, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0b, 0x0a, 0x07, 0x46, 0x55, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07,
	0x53, 0x54, 0x41, 0x4b, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x44, 0x45, 0x58,
	0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x47, 0x4f, 0x56, 0x45, 0x52, 0x4e, 0x41, 0x4e, 0x43, 0x45,
	0x10, 0x03, 0x42, 0xaf, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73,
	0x63, 0x6f, 0x70, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x13, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x6f,
	0x6f, 0x6d, 0x69, 0x2d, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x72, 0x2d, 0x73, 0x63,
	0x6f, 0x70, 0x65, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0xa2, 0x02, 0x03, 0x53, 0x45, 0x58, 0xaa,
	0x02, 0x0f, 0x53, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0xca, 0x02, 0x0f, 0x53, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x5c, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0xe2, 0x02, 0x1b, 0x53, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x5c,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x10, 0x53, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x3a, 0x3a, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_event_processed_event_proto_rawDescOnce sync.Once
	file_event_processed_event_proto_rawDescData = file_event_processed_event_proto_rawDesc
)

func file_event_processed_event_proto_rawDescGZIP() []byte {
	file_event_processed_event_proto_rawDescOnce.Do(func() {
		file_event_processed_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_event_processed_event_proto_rawDescData)
	})
	return file_event_processed_event_proto_rawDescData
}

var file_event_processed_event_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_event_processed_event_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_event_processed_event_proto_goTypes = []interface{}{
	(EventType)(0),            // 0: starscope.event.EventType
	(*EventProcessedMsg)(nil), // 1: starscope.event.EventProcessedMsg
}
var file_event_processed_event_proto_depIdxs = []int32{
	0, // 0: starscope.event.EventProcessedMsg.event_type:type_name -> starscope.event.EventType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_event_processed_event_proto_init() }
func file_event_processed_event_proto_init() {
	if File_event_processed_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_event_processed_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventProcessedMsg); i {
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
			RawDescriptor: file_event_processed_event_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_event_processed_event_proto_goTypes,
		DependencyIndexes: file_event_processed_event_proto_depIdxs,
		EnumInfos:         file_event_processed_event_proto_enumTypes,
		MessageInfos:      file_event_processed_event_proto_msgTypes,
	}.Build()
	File_event_processed_event_proto = out.File
	file_event_processed_event_proto_rawDesc = nil
	file_event_processed_event_proto_goTypes = nil
	file_event_processed_event_proto_depIdxs = nil
}
