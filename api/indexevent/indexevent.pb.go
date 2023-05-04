// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: indexevent/indexevent.proto

package indexevent

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

type Coin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Denom  string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Amount string `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *Coin) Reset() {
	*x = Coin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indexevent_indexevent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Coin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Coin) ProtoMessage() {}

func (x *Coin) ProtoReflect() protoreflect.Message {
	mi := &file_indexevent_indexevent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Coin.ProtoReflect.Descriptor instead.
func (*Coin) Descriptor() ([]byte, []int) {
	return file_indexevent_indexevent_proto_rawDescGZIP(), []int{0}
}

func (x *Coin) GetDenom() string {
	if x != nil {
		return x.Denom
	}
	return ""
}

func (x *Coin) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

type CoinReceivedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Coin   *Coin  `protobuf:"bytes,2,opt,name=coin,proto3" json:"coin,omitempty"`
}

func (x *CoinReceivedEvent) Reset() {
	*x = CoinReceivedEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indexevent_indexevent_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CoinReceivedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CoinReceivedEvent) ProtoMessage() {}

func (x *CoinReceivedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_indexevent_indexevent_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CoinReceivedEvent.ProtoReflect.Descriptor instead.
func (*CoinReceivedEvent) Descriptor() ([]byte, []int) {
	return file_indexevent_indexevent_proto_rawDescGZIP(), []int{1}
}

func (x *CoinReceivedEvent) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *CoinReceivedEvent) GetCoin() *Coin {
	if x != nil {
		return x.Coin
	}
	return nil
}

type OsmosisPoolUnlockEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Duration   *durationpb.Duration   `protobuf:"bytes,1,opt,name=duration,proto3" json:"duration,omitempty"`
	UnlockTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=unlock_time,json=unlockTime,proto3" json:"unlock_time,omitempty"`
}

func (x *OsmosisPoolUnlockEvent) Reset() {
	*x = OsmosisPoolUnlockEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indexevent_indexevent_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OsmosisPoolUnlockEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OsmosisPoolUnlockEvent) ProtoMessage() {}

func (x *OsmosisPoolUnlockEvent) ProtoReflect() protoreflect.Message {
	mi := &file_indexevent_indexevent_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OsmosisPoolUnlockEvent.ProtoReflect.Descriptor instead.
func (*OsmosisPoolUnlockEvent) Descriptor() ([]byte, []int) {
	return file_indexevent_indexevent_proto_rawDescGZIP(), []int{2}
}

func (x *OsmosisPoolUnlockEvent) GetDuration() *durationpb.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *OsmosisPoolUnlockEvent) GetUnlockTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UnlockTime
	}
	return nil
}

type TxEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChainName     string                 `protobuf:"bytes,1,opt,name=chain_name,json=chainName,proto3" json:"chain_name,omitempty"`
	WalletAddress string                 `protobuf:"bytes,2,opt,name=wallet_address,json=walletAddress,proto3" json:"wallet_address,omitempty"`
	NotifyTime    *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=notify_time,json=notifyTime,proto3" json:"notify_time,omitempty"`
	// Types that are assignable to Event:
	//
	//	*TxEvent_CoinReceived
	//	*TxEvent_OsmosisPoolUnlock
	Event isTxEvent_Event `protobuf_oneof:"event"`
}

func (x *TxEvent) Reset() {
	*x = TxEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indexevent_indexevent_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TxEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TxEvent) ProtoMessage() {}

func (x *TxEvent) ProtoReflect() protoreflect.Message {
	mi := &file_indexevent_indexevent_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TxEvent.ProtoReflect.Descriptor instead.
func (*TxEvent) Descriptor() ([]byte, []int) {
	return file_indexevent_indexevent_proto_rawDescGZIP(), []int{3}
}

func (x *TxEvent) GetChainName() string {
	if x != nil {
		return x.ChainName
	}
	return ""
}

func (x *TxEvent) GetWalletAddress() string {
	if x != nil {
		return x.WalletAddress
	}
	return ""
}

func (x *TxEvent) GetNotifyTime() *timestamppb.Timestamp {
	if x != nil {
		return x.NotifyTime
	}
	return nil
}

func (m *TxEvent) GetEvent() isTxEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (x *TxEvent) GetCoinReceived() *CoinReceivedEvent {
	if x, ok := x.GetEvent().(*TxEvent_CoinReceived); ok {
		return x.CoinReceived
	}
	return nil
}

func (x *TxEvent) GetOsmosisPoolUnlock() *OsmosisPoolUnlockEvent {
	if x, ok := x.GetEvent().(*TxEvent_OsmosisPoolUnlock); ok {
		return x.OsmosisPoolUnlock
	}
	return nil
}

type isTxEvent_Event interface {
	isTxEvent_Event()
}

type TxEvent_CoinReceived struct {
	CoinReceived *CoinReceivedEvent `protobuf:"bytes,4,opt,name=coin_received,json=coinReceived,proto3,oneof"`
}

type TxEvent_OsmosisPoolUnlock struct {
	OsmosisPoolUnlock *OsmosisPoolUnlockEvent `protobuf:"bytes,5,opt,name=osmosis_pool_unlock,json=osmosisPoolUnlock,proto3,oneof"`
}

func (*TxEvent_CoinReceived) isTxEvent_Event() {}

func (*TxEvent_OsmosisPoolUnlock) isTxEvent_Event() {}

var File_indexevent_indexevent_proto protoreflect.FileDescriptor

var file_indexevent_indexevent_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2f, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x73,
	0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x34, 0x0a, 0x04, 0x43, 0x6f, 0x69, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6e, 0x6f, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x12, 0x16, 0x0a,
	0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x56, 0x0a, 0x11, 0x43, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x12, 0x29, 0x0a, 0x04, 0x63, 0x6f, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x15, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x69, 0x6e, 0x52, 0x04, 0x63, 0x6f, 0x69, 0x6e, 0x22, 0x8c, 0x01,
	0x0a, 0x16, 0x4f, 0x73, 0x6d, 0x6f, 0x73, 0x69, 0x73, 0x50, 0x6f, 0x6f, 0x6c, 0x55, 0x6e, 0x6c,
	0x6f, 0x63, 0x6b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x35, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x3b, 0x0a, 0x0b, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x0a, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x22, 0xbb, 0x02, 0x0a,
	0x07, 0x54, 0x78, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x3b,
	0x0a, 0x0b, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0a, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x49, 0x0a, 0x0d, 0x63,
	0x6f, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x22, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x0c, 0x63, 0x6f, 0x69, 0x6e, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x12, 0x59, 0x0a, 0x13, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x69,
	0x73, 0x5f, 0x70, 0x6f, 0x6f, 0x6c, 0x5f, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x4f, 0x73, 0x6d, 0x6f, 0x73, 0x69, 0x73, 0x50, 0x6f, 0x6f,
	0x6c, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x11,
	0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x69, 0x73, 0x50, 0x6f, 0x6f, 0x6c, 0x55, 0x6e, 0x6c, 0x6f, 0x63,
	0x6b, 0x42, 0x07, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x42, 0xb0, 0x01, 0x0a, 0x13, 0x63,
	0x6f, 0x6d, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x42, 0x0f, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6c, 0x6f, 0x6f, 0x6d, 0x69, 0x2d, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x73, 0x74, 0x61,
	0x72, 0x2d, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0xa2, 0x02, 0x03, 0x53, 0x45, 0x58, 0xaa, 0x02, 0x0f, 0x53, 0x74, 0x61, 0x72, 0x73,
	0x63, 0x6f, 0x70, 0x65, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0xca, 0x02, 0x0f, 0x53, 0x74, 0x61,
	0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x5c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0xe2, 0x02, 0x1b, 0x53,
	0x74, 0x61, 0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x5c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x5c, 0x47,
	0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x10, 0x53, 0x74, 0x61,
	0x72, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x3a, 0x3a, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_indexevent_indexevent_proto_rawDescOnce sync.Once
	file_indexevent_indexevent_proto_rawDescData = file_indexevent_indexevent_proto_rawDesc
)

func file_indexevent_indexevent_proto_rawDescGZIP() []byte {
	file_indexevent_indexevent_proto_rawDescOnce.Do(func() {
		file_indexevent_indexevent_proto_rawDescData = protoimpl.X.CompressGZIP(file_indexevent_indexevent_proto_rawDescData)
	})
	return file_indexevent_indexevent_proto_rawDescData
}

var file_indexevent_indexevent_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_indexevent_indexevent_proto_goTypes = []interface{}{
	(*Coin)(nil),                   // 0: starscope.event.Coin
	(*CoinReceivedEvent)(nil),      // 1: starscope.event.CoinReceivedEvent
	(*OsmosisPoolUnlockEvent)(nil), // 2: starscope.event.OsmosisPoolUnlockEvent
	(*TxEvent)(nil),                // 3: starscope.event.TxEvent
	(*durationpb.Duration)(nil),    // 4: google.protobuf.Duration
	(*timestamppb.Timestamp)(nil),  // 5: google.protobuf.Timestamp
}
var file_indexevent_indexevent_proto_depIdxs = []int32{
	0, // 0: starscope.event.CoinReceivedEvent.coin:type_name -> starscope.event.Coin
	4, // 1: starscope.event.OsmosisPoolUnlockEvent.duration:type_name -> google.protobuf.Duration
	5, // 2: starscope.event.OsmosisPoolUnlockEvent.unlock_time:type_name -> google.protobuf.Timestamp
	5, // 3: starscope.event.TxEvent.notify_time:type_name -> google.protobuf.Timestamp
	1, // 4: starscope.event.TxEvent.coin_received:type_name -> starscope.event.CoinReceivedEvent
	2, // 5: starscope.event.TxEvent.osmosis_pool_unlock:type_name -> starscope.event.OsmosisPoolUnlockEvent
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_indexevent_indexevent_proto_init() }
func file_indexevent_indexevent_proto_init() {
	if File_indexevent_indexevent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_indexevent_indexevent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Coin); i {
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
		file_indexevent_indexevent_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CoinReceivedEvent); i {
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
		file_indexevent_indexevent_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OsmosisPoolUnlockEvent); i {
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
		file_indexevent_indexevent_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TxEvent); i {
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
	file_indexevent_indexevent_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*TxEvent_CoinReceived)(nil),
		(*TxEvent_OsmosisPoolUnlock)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_indexevent_indexevent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_indexevent_indexevent_proto_goTypes,
		DependencyIndexes: file_indexevent_indexevent_proto_depIdxs,
		MessageInfos:      file_indexevent_indexevent_proto_msgTypes,
	}.Build()
	File_indexevent_indexevent_proto = out.File
	file_indexevent_indexevent_proto_rawDesc = nil
	file_indexevent_indexevent_proto_goTypes = nil
	file_indexevent_indexevent_proto_depIdxs = nil
}
