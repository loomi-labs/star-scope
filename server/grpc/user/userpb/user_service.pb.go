// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: grpc/user/userpb/user_service.proto

package userpb

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

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name            string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	HasDiscord      bool   `protobuf:"varint,3,opt,name=has_discord,json=hasDiscord,proto3" json:"has_discord,omitempty"`
	HasTelegram     bool   `protobuf:"varint,4,opt,name=has_telegram,json=hasTelegram,proto3" json:"has_telegram,omitempty"`
	IsSetupComplete bool   `protobuf:"varint,5,opt,name=is_setup_complete,json=isSetupComplete,proto3" json:"is_setup_complete,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_user_userpb_user_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_user_userpb_user_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_grpc_user_userpb_user_service_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetHasDiscord() bool {
	if x != nil {
		return x.HasDiscord
	}
	return false
}

func (x *User) GetHasTelegram() bool {
	if x != nil {
		return x.HasTelegram
	}
	return false
}

func (x *User) GetIsSetupComplete() bool {
	if x != nil {
		return x.IsSetupComplete
	}
	return false
}

type DiscordChannel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ChannelId int64  `protobuf:"varint,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	Name      string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	IsGroup   bool   `protobuf:"varint,4,opt,name=is_group,json=isGroup,proto3" json:"is_group,omitempty"`
}

func (x *DiscordChannel) Reset() {
	*x = DiscordChannel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_user_userpb_user_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscordChannel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscordChannel) ProtoMessage() {}

func (x *DiscordChannel) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_user_userpb_user_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscordChannel.ProtoReflect.Descriptor instead.
func (*DiscordChannel) Descriptor() ([]byte, []int) {
	return file_grpc_user_userpb_user_service_proto_rawDescGZIP(), []int{1}
}

func (x *DiscordChannel) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DiscordChannel) GetChannelId() int64 {
	if x != nil {
		return x.ChannelId
	}
	return 0
}

func (x *DiscordChannel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DiscordChannel) GetIsGroup() bool {
	if x != nil {
		return x.IsGroup
	}
	return false
}

type DiscordChannels struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Channels []*DiscordChannel `protobuf:"bytes,1,rep,name=channels,proto3" json:"channels,omitempty"`
}

func (x *DiscordChannels) Reset() {
	*x = DiscordChannels{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_user_userpb_user_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscordChannels) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscordChannels) ProtoMessage() {}

func (x *DiscordChannels) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_user_userpb_user_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscordChannels.ProtoReflect.Descriptor instead.
func (*DiscordChannels) Descriptor() ([]byte, []int) {
	return file_grpc_user_userpb_user_service_proto_rawDescGZIP(), []int{2}
}

func (x *DiscordChannels) GetChannels() []*DiscordChannel {
	if x != nil {
		return x.Channels
	}
	return nil
}

type DeleteDiscordChannelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId int64 `protobuf:"varint,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
}

func (x *DeleteDiscordChannelRequest) Reset() {
	*x = DeleteDiscordChannelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_user_userpb_user_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteDiscordChannelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteDiscordChannelRequest) ProtoMessage() {}

func (x *DeleteDiscordChannelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_user_userpb_user_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteDiscordChannelRequest.ProtoReflect.Descriptor instead.
func (*DeleteDiscordChannelRequest) Descriptor() ([]byte, []int) {
	return file_grpc_user_userpb_user_service_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteDiscordChannelRequest) GetChannelId() int64 {
	if x != nil {
		return x.ChannelId
	}
	return 0
}

type TelegramChat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ChatId  int64  `protobuf:"varint,2,opt,name=chat_id,json=chatId,proto3" json:"chat_id,omitempty"`
	Name    string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	IsGroup bool   `protobuf:"varint,4,opt,name=is_group,json=isGroup,proto3" json:"is_group,omitempty"`
}

func (x *TelegramChat) Reset() {
	*x = TelegramChat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_user_userpb_user_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelegramChat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelegramChat) ProtoMessage() {}

func (x *TelegramChat) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_user_userpb_user_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelegramChat.ProtoReflect.Descriptor instead.
func (*TelegramChat) Descriptor() ([]byte, []int) {
	return file_grpc_user_userpb_user_service_proto_rawDescGZIP(), []int{4}
}

func (x *TelegramChat) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *TelegramChat) GetChatId() int64 {
	if x != nil {
		return x.ChatId
	}
	return 0
}

func (x *TelegramChat) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TelegramChat) GetIsGroup() bool {
	if x != nil {
		return x.IsGroup
	}
	return false
}

type TelegramChats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chats []*TelegramChat `protobuf:"bytes,1,rep,name=chats,proto3" json:"chats,omitempty"`
}

func (x *TelegramChats) Reset() {
	*x = TelegramChats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_user_userpb_user_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelegramChats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelegramChats) ProtoMessage() {}

func (x *TelegramChats) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_user_userpb_user_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelegramChats.ProtoReflect.Descriptor instead.
func (*TelegramChats) Descriptor() ([]byte, []int) {
	return file_grpc_user_userpb_user_service_proto_rawDescGZIP(), []int{5}
}

func (x *TelegramChats) GetChats() []*TelegramChat {
	if x != nil {
		return x.Chats
	}
	return nil
}

type DeleteTelegramChatRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChatId int64 `protobuf:"varint,1,opt,name=chat_id,json=chatId,proto3" json:"chat_id,omitempty"`
}

func (x *DeleteTelegramChatRequest) Reset() {
	*x = DeleteTelegramChatRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_user_userpb_user_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteTelegramChatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTelegramChatRequest) ProtoMessage() {}

func (x *DeleteTelegramChatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_user_userpb_user_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTelegramChatRequest.ProtoReflect.Descriptor instead.
func (*DeleteTelegramChatRequest) Descriptor() ([]byte, []int) {
	return file_grpc_user_userpb_user_service_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteTelegramChatRequest) GetChatId() int64 {
	if x != nil {
		return x.ChatId
	}
	return 0
}

var File_grpc_user_userpb_user_service_proto protoreflect.FileDescriptor

var file_grpc_user_userpb_user_service_proto_rawDesc = []byte{
	0x0a, 0x23, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x70, 0x62, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x9a, 0x01, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x68, 0x61, 0x73, 0x5f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x68, 0x61, 0x73, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64,
	0x12, 0x21, 0x0a, 0x0c, 0x68, 0x61, 0x73, 0x5f, 0x74, 0x65, 0x6c, 0x65, 0x67, 0x72, 0x61, 0x6d,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x68, 0x61, 0x73, 0x54, 0x65, 0x6c, 0x65, 0x67,
	0x72, 0x61, 0x6d, 0x12, 0x2a, 0x0a, 0x11, 0x69, 0x73, 0x5f, 0x73, 0x65, 0x74, 0x75, 0x70, 0x5f,
	0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f,
	0x69, 0x73, 0x53, 0x65, 0x74, 0x75, 0x70, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x22,
	0x6e, 0x0a, 0x0e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22,
	0x4d, 0x0a, 0x0f, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x73, 0x12, 0x3a, 0x0a, 0x08, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x08, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x22, 0x3c,
	0x0a, 0x1b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x43,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64, 0x22, 0x66, 0x0a, 0x0c,
	0x54, 0x65, 0x6c, 0x65, 0x67, 0x72, 0x61, 0x6d, 0x43, 0x68, 0x61, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07,
	0x63, 0x68, 0x61, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x63,
	0x68, 0x61, 0x74, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x22, 0x43, 0x0a, 0x0d, 0x54, 0x65, 0x6c, 0x65, 0x67, 0x72, 0x61, 0x6d,
	0x43, 0x68, 0x61, 0x74, 0x73, 0x12, 0x32, 0x0a, 0x05, 0x63, 0x68, 0x61, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x65, 0x6c, 0x65, 0x67, 0x72, 0x61, 0x6d, 0x43, 0x68,
	0x61, 0x74, 0x52, 0x05, 0x63, 0x68, 0x61, 0x74, 0x73, 0x22, 0x34, 0x0a, 0x19, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x54, 0x65, 0x6c, 0x65, 0x67, 0x72, 0x61, 0x6d, 0x43, 0x68, 0x61, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x32,
	0xe5, 0x03, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x39, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x14, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x0d, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x50, 0x0a,
	0x13, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1f, 0x2e, 0x73,
	0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x22, 0x00, 0x12,
	0x5d, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64,
	0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x2b, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63,
	0x6f, 0x70, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x4c,
	0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x6c, 0x65, 0x67, 0x72, 0x61, 0x6d, 0x43, 0x68,
	0x61, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1d, 0x2e, 0x73, 0x74,
	0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x65, 0x6c,
	0x65, 0x67, 0x72, 0x61, 0x6d, 0x43, 0x68, 0x61, 0x74, 0x73, 0x22, 0x00, 0x12, 0x59, 0x0a, 0x12,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x6c, 0x65, 0x67, 0x72, 0x61, 0x6d, 0x43, 0x68,
	0x61, 0x74, 0x12, 0x29, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x6c, 0x65, 0x67, 0x72,
	0x61, 0x6d, 0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0xb2, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e,
	0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x42, 0x10,
	0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c,
	0x6f, 0x6f, 0x6d, 0x69, 0x2d, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x72, 0x2d, 0x73,
	0x63, 0x6f, 0x70, 0x65, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x70, 0x62, 0xa2, 0x02, 0x03, 0x53, 0x47, 0x58, 0xaa, 0x02, 0x0e, 0x53, 0x74,
	0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x47, 0x72, 0x70, 0x63, 0xca, 0x02, 0x0e, 0x53,
	0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x5c, 0x47, 0x72, 0x70, 0x63, 0xe2, 0x02, 0x1a,
	0x53, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x5c, 0x47, 0x72, 0x70, 0x63, 0x5c, 0x47,
	0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0f, 0x53, 0x74, 0x61,
	0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x3a, 0x3a, 0x47, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_user_userpb_user_service_proto_rawDescOnce sync.Once
	file_grpc_user_userpb_user_service_proto_rawDescData = file_grpc_user_userpb_user_service_proto_rawDesc
)

func file_grpc_user_userpb_user_service_proto_rawDescGZIP() []byte {
	file_grpc_user_userpb_user_service_proto_rawDescOnce.Do(func() {
		file_grpc_user_userpb_user_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_user_userpb_user_service_proto_rawDescData)
	})
	return file_grpc_user_userpb_user_service_proto_rawDescData
}

var file_grpc_user_userpb_user_service_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_grpc_user_userpb_user_service_proto_goTypes = []interface{}{
	(*User)(nil),                        // 0: starscope.grpc.User
	(*DiscordChannel)(nil),              // 1: starscope.grpc.DiscordChannel
	(*DiscordChannels)(nil),             // 2: starscope.grpc.DiscordChannels
	(*DeleteDiscordChannelRequest)(nil), // 3: starscope.grpc.DeleteDiscordChannelRequest
	(*TelegramChat)(nil),                // 4: starscope.grpc.TelegramChat
	(*TelegramChats)(nil),               // 5: starscope.grpc.TelegramChats
	(*DeleteTelegramChatRequest)(nil),   // 6: starscope.grpc.DeleteTelegramChatRequest
	(*emptypb.Empty)(nil),               // 7: google.protobuf.Empty
}
var file_grpc_user_userpb_user_service_proto_depIdxs = []int32{
	1, // 0: starscope.grpc.DiscordChannels.channels:type_name -> starscope.grpc.DiscordChannel
	4, // 1: starscope.grpc.TelegramChats.chats:type_name -> starscope.grpc.TelegramChat
	7, // 2: starscope.grpc.UserService.GetUser:input_type -> google.protobuf.Empty
	7, // 3: starscope.grpc.UserService.DeleteAccount:input_type -> google.protobuf.Empty
	7, // 4: starscope.grpc.UserService.ListDiscordChannels:input_type -> google.protobuf.Empty
	3, // 5: starscope.grpc.UserService.DeleteDiscordChannel:input_type -> starscope.grpc.DeleteDiscordChannelRequest
	7, // 6: starscope.grpc.UserService.ListTelegramChats:input_type -> google.protobuf.Empty
	6, // 7: starscope.grpc.UserService.DeleteTelegramChat:input_type -> starscope.grpc.DeleteTelegramChatRequest
	0, // 8: starscope.grpc.UserService.GetUser:output_type -> starscope.grpc.User
	7, // 9: starscope.grpc.UserService.DeleteAccount:output_type -> google.protobuf.Empty
	2, // 10: starscope.grpc.UserService.ListDiscordChannels:output_type -> starscope.grpc.DiscordChannels
	7, // 11: starscope.grpc.UserService.DeleteDiscordChannel:output_type -> google.protobuf.Empty
	5, // 12: starscope.grpc.UserService.ListTelegramChats:output_type -> starscope.grpc.TelegramChats
	7, // 13: starscope.grpc.UserService.DeleteTelegramChat:output_type -> google.protobuf.Empty
	8, // [8:14] is the sub-list for method output_type
	2, // [2:8] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_grpc_user_userpb_user_service_proto_init() }
func file_grpc_user_userpb_user_service_proto_init() {
	if File_grpc_user_userpb_user_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_user_userpb_user_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_grpc_user_userpb_user_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscordChannel); i {
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
		file_grpc_user_userpb_user_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscordChannels); i {
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
		file_grpc_user_userpb_user_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteDiscordChannelRequest); i {
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
		file_grpc_user_userpb_user_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelegramChat); i {
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
		file_grpc_user_userpb_user_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelegramChats); i {
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
		file_grpc_user_userpb_user_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteTelegramChatRequest); i {
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
			RawDescriptor: file_grpc_user_userpb_user_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_user_userpb_user_service_proto_goTypes,
		DependencyIndexes: file_grpc_user_userpb_user_service_proto_depIdxs,
		MessageInfos:      file_grpc_user_userpb_user_service_proto_msgTypes,
	}.Build()
	File_grpc_user_userpb_user_service_proto = out.File
	file_grpc_user_userpb_user_service_proto_rawDesc = nil
	file_grpc_user_userpb_user_service_proto_goTypes = nil
	file_grpc_user_userpb_user_service_proto_depIdxs = nil
}
