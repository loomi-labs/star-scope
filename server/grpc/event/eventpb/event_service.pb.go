// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: grpc/event/eventpb/event_service.proto

package eventpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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
	return file_grpc_event_eventpb_event_service_proto_enumTypes[0].Descriptor()
}

func (EventType) Type() protoreflect.EnumType {
	return &file_grpc_event_eventpb_event_service_proto_enumTypes[0]
}

func (x EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventType.Descriptor instead.
func (EventType) EnumDescriptor() ([]byte, []int) {
	return file_grpc_event_eventpb_event_service_proto_rawDescGZIP(), []int{0}
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title       string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Subtitle    string                 `protobuf:"bytes,3,opt,name=subtitle,proto3" json:"subtitle,omitempty"`
	Description string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Timestamp   *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	EventType   EventType              `protobuf:"varint,6,opt,name=event_type,json=eventType,proto3,enum=starscope.grpc.EventType" json:"event_type,omitempty"`
	Chain       *ChainData             `protobuf:"bytes,7,opt,name=chain,proto3" json:"chain,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_event_eventpb_event_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_event_eventpb_event_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_grpc_event_eventpb_event_service_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Event) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Event) GetSubtitle() string {
	if x != nil {
		return x.Subtitle
	}
	return ""
}

func (x *Event) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Event) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Event) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_FUNDING
}

func (x *Event) GetChain() *ChainData {
	if x != nil {
		return x.Chain
	}
	return nil
}

type EventList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Events []*Event `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
}

func (x *EventList) Reset() {
	*x = EventList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_event_eventpb_event_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventList) ProtoMessage() {}

func (x *EventList) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_event_eventpb_event_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventList.ProtoReflect.Descriptor instead.
func (*EventList) Descriptor() ([]byte, []int) {
	return file_grpc_event_eventpb_event_service_proto_rawDescGZIP(), []int{1}
}

func (x *EventList) GetEvents() []*Event {
	if x != nil {
		return x.Events
	}
	return nil
}

type ListEventsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
}

func (x *ListEventsRequest) Reset() {
	*x = ListEventsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_event_eventpb_event_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEventsRequest) ProtoMessage() {}

func (x *ListEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_event_eventpb_event_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListEventsRequest.ProtoReflect.Descriptor instead.
func (*ListEventsRequest) Descriptor() ([]byte, []int) {
	return file_grpc_event_eventpb_event_service_proto_rawDescGZIP(), []int{2}
}

func (x *ListEventsRequest) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

type ChainData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	ImageUrl string `protobuf:"bytes,3,opt,name=image_url,json=imageUrl,proto3" json:"image_url,omitempty"`
}

func (x *ChainData) Reset() {
	*x = ChainData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_event_eventpb_event_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChainData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChainData) ProtoMessage() {}

func (x *ChainData) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_event_eventpb_event_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChainData.ProtoReflect.Descriptor instead.
func (*ChainData) Descriptor() ([]byte, []int) {
	return file_grpc_event_eventpb_event_service_proto_rawDescGZIP(), []int{3}
}

func (x *ChainData) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ChainData) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChainData) GetImageUrl() string {
	if x != nil {
		return x.ImageUrl
	}
	return ""
}

type ChainList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chains []*ChainData `protobuf:"bytes,1,rep,name=chains,proto3" json:"chains,omitempty"`
}

func (x *ChainList) Reset() {
	*x = ChainList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_event_eventpb_event_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChainList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChainList) ProtoMessage() {}

func (x *ChainList) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_event_eventpb_event_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChainList.ProtoReflect.Descriptor instead.
func (*ChainList) Descriptor() ([]byte, []int) {
	return file_grpc_event_eventpb_event_service_proto_rawDescGZIP(), []int{4}
}

func (x *ChainList) GetChains() []*ChainData {
	if x != nil {
		return x.Chains
	}
	return nil
}

var File_grpc_event_eventpb_event_service_proto protoreflect.FileDescriptor

var file_grpc_event_eventpb_event_service_proto_rawDesc = []byte{
	0x0a, 0x26, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x70, 0x62, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63,
	0x6f, 0x70, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x90, 0x02, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x75, 0x62, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x75, 0x62, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x38,
	0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x19, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2f, 0x0a, 0x05, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63,
	0x6f, 0x70, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x44, 0x61,
	0x74, 0x61, 0x52, 0x05, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x22, 0x3a, 0x0a, 0x09, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f,
	0x70, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x06, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x4e, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x4c, 0x0a, 0x09, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f,
	0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65,
	0x55, 0x72, 0x6c, 0x22, 0x3e, 0x0a, 0x09, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x31, 0x0a, 0x06, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x06, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x73, 0x2a, 0x3e, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x0b, 0x0a, 0x07, 0x46, 0x55, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x00, 0x12, 0x0b, 0x0a,
	0x07, 0x53, 0x54, 0x41, 0x4b, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x44, 0x45,
	0x58, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x47, 0x4f, 0x56, 0x45, 0x52, 0x4e, 0x41, 0x4e, 0x43,
	0x45, 0x10, 0x03, 0x32, 0xe5, 0x01, 0x0a, 0x0c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x44, 0x0a, 0x0b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e, 0x73, 0x74,
	0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x30, 0x01, 0x12, 0x4c, 0x0a, 0x0a, 0x4c, 0x69,
	0x73, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x21, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73,
	0x63, 0x6f, 0x70, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x73, 0x74,
	0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74,
	0x43, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19,
	0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x43, 0x68, 0x61, 0x69, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x42, 0xb5, 0x01, 0x0a, 0x12,
	0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x42, 0x11, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x6f, 0x6f, 0x6d, 0x69, 0x2d, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x73,
	0x74, 0x61, 0x72, 0x2d, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70, 0x62, 0xa2, 0x02, 0x03, 0x53,
	0x47, 0x58, 0xaa, 0x02, 0x0e, 0x53, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x47,
	0x72, 0x70, 0x63, 0xca, 0x02, 0x0e, 0x53, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x5c,
	0x47, 0x72, 0x70, 0x63, 0xe2, 0x02, 0x1a, 0x53, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65,
	0x5c, 0x47, 0x72, 0x70, 0x63, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x0f, 0x53, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x3a, 0x3a, 0x47,
	0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_event_eventpb_event_service_proto_rawDescOnce sync.Once
	file_grpc_event_eventpb_event_service_proto_rawDescData = file_grpc_event_eventpb_event_service_proto_rawDesc
)

func file_grpc_event_eventpb_event_service_proto_rawDescGZIP() []byte {
	file_grpc_event_eventpb_event_service_proto_rawDescOnce.Do(func() {
		file_grpc_event_eventpb_event_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_event_eventpb_event_service_proto_rawDescData)
	})
	return file_grpc_event_eventpb_event_service_proto_rawDescData
}

var file_grpc_event_eventpb_event_service_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_grpc_event_eventpb_event_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_grpc_event_eventpb_event_service_proto_goTypes = []interface{}{
	(EventType)(0),                // 0: starscope.grpc.EventType
	(*Event)(nil),                 // 1: starscope.grpc.Event
	(*EventList)(nil),             // 2: starscope.grpc.EventList
	(*ListEventsRequest)(nil),     // 3: starscope.grpc.ListEventsRequest
	(*ChainData)(nil),             // 4: starscope.grpc.ChainData
	(*ChainList)(nil),             // 5: starscope.grpc.ChainList
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 7: google.protobuf.Empty
}
var file_grpc_event_eventpb_event_service_proto_depIdxs = []int32{
	6, // 0: starscope.grpc.Event.timestamp:type_name -> google.protobuf.Timestamp
	0, // 1: starscope.grpc.Event.event_type:type_name -> starscope.grpc.EventType
	4, // 2: starscope.grpc.Event.chain:type_name -> starscope.grpc.ChainData
	1, // 3: starscope.grpc.EventList.events:type_name -> starscope.grpc.Event
	6, // 4: starscope.grpc.ListEventsRequest.start_time:type_name -> google.protobuf.Timestamp
	4, // 5: starscope.grpc.ChainList.chains:type_name -> starscope.grpc.ChainData
	7, // 6: starscope.grpc.EventService.EventStream:input_type -> google.protobuf.Empty
	3, // 7: starscope.grpc.EventService.ListEvents:input_type -> starscope.grpc.ListEventsRequest
	7, // 8: starscope.grpc.EventService.ListChains:input_type -> google.protobuf.Empty
	2, // 9: starscope.grpc.EventService.EventStream:output_type -> starscope.grpc.EventList
	2, // 10: starscope.grpc.EventService.ListEvents:output_type -> starscope.grpc.EventList
	5, // 11: starscope.grpc.EventService.ListChains:output_type -> starscope.grpc.ChainList
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_grpc_event_eventpb_event_service_proto_init() }
func file_grpc_event_eventpb_event_service_proto_init() {
	if File_grpc_event_eventpb_event_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_event_eventpb_event_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
		file_grpc_event_eventpb_event_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventList); i {
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
		file_grpc_event_eventpb_event_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListEventsRequest); i {
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
		file_grpc_event_eventpb_event_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChainData); i {
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
		file_grpc_event_eventpb_event_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChainList); i {
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
			RawDescriptor: file_grpc_event_eventpb_event_service_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_event_eventpb_event_service_proto_goTypes,
		DependencyIndexes: file_grpc_event_eventpb_event_service_proto_depIdxs,
		EnumInfos:         file_grpc_event_eventpb_event_service_proto_enumTypes,
		MessageInfos:      file_grpc_event_eventpb_event_service_proto_msgTypes,
	}.Build()
	File_grpc_event_eventpb_event_service_proto = out.File
	file_grpc_event_eventpb_event_service_proto_rawDesc = nil
	file_grpc_event_eventpb_event_service_proto_goTypes = nil
	file_grpc_event_eventpb_event_service_proto_depIdxs = nil
}
