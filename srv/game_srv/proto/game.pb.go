// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: game.proto

package proto

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

type Type int32

const (
	Type_Apple  Type = 0
	Type_Banana Type = 1
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "Apple",
		1: "Banana",
	}
	Type_value = map[string]int32{
		"Apple":  0,
		"Banana": 1,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_game_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_game_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{0}
}

type UserItemsInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Gold    uint32 `protobuf:"varint,2,opt,name=gold,proto3" json:"gold,omitempty"`
	Diamond uint32 `protobuf:"varint,3,opt,name=diamond,proto3" json:"diamond,omitempty"`
	Apple   uint32 `protobuf:"varint,4,opt,name=apple,proto3" json:"apple,omitempty"`
	Banana  uint32 `protobuf:"varint,5,opt,name=banana,proto3" json:"banana,omitempty"`
}

func (x *UserItemsInfo) Reset() {
	*x = UserItemsInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserItemsInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserItemsInfo) ProtoMessage() {}

func (x *UserItemsInfo) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserItemsInfo.ProtoReflect.Descriptor instead.
func (*UserItemsInfo) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{0}
}

func (x *UserItemsInfo) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserItemsInfo) GetGold() uint32 {
	if x != nil {
		return x.Gold
	}
	return 0
}

func (x *UserItemsInfo) GetDiamond() uint32 {
	if x != nil {
		return x.Diamond
	}
	return 0
}

func (x *UserItemsInfo) GetApple() uint32 {
	if x != nil {
		return x.Apple
	}
	return 0
}

func (x *UserItemsInfo) GetBanana() uint32 {
	if x != nil {
		return x.Banana
	}
	return 0
}

type UserIDInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *UserIDInfo) Reset() {
	*x = UserIDInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserIDInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserIDInfo) ProtoMessage() {}

func (x *UserIDInfo) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserIDInfo.ProtoReflect.Descriptor instead.
func (*UserIDInfo) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{1}
}

func (x *UserIDInfo) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type UserItemsInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Gold    uint32 `protobuf:"varint,2,opt,name=gold,proto3" json:"gold,omitempty"`
	Diamond uint32 `protobuf:"varint,3,opt,name=diamond,proto3" json:"diamond,omitempty"`
	// 道具
	Items []uint32 `protobuf:"varint,4,rep,packed,name=items,proto3" json:"items,omitempty"`
}

func (x *UserItemsInfoResponse) Reset() {
	*x = UserItemsInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserItemsInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserItemsInfoResponse) ProtoMessage() {}

func (x *UserItemsInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserItemsInfoResponse.ProtoReflect.Descriptor instead.
func (*UserItemsInfoResponse) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{2}
}

func (x *UserItemsInfoResponse) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserItemsInfoResponse) GetGold() uint32 {
	if x != nil {
		return x.Gold
	}
	return 0
}

func (x *UserItemsInfoResponse) GetDiamond() uint32 {
	if x != nil {
		return x.Diamond
	}
	return 0
}

func (x *UserItemsInfoResponse) GetItems() []uint32 {
	if x != nil {
		return x.Items
	}
	return nil
}

type AddGoldRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Count uint32 `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *AddGoldRequest) Reset() {
	*x = AddGoldRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddGoldRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddGoldRequest) ProtoMessage() {}

func (x *AddGoldRequest) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddGoldRequest.ProtoReflect.Descriptor instead.
func (*AddGoldRequest) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{3}
}

func (x *AddGoldRequest) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AddGoldRequest) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type AddDiamondInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Count uint32 `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *AddDiamondInfo) Reset() {
	*x = AddDiamondInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddDiamondInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddDiamondInfo) ProtoMessage() {}

func (x *AddDiamondInfo) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddDiamondInfo.ProtoReflect.Descriptor instead.
func (*AddDiamondInfo) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{4}
}

func (x *AddDiamondInfo) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AddDiamondInfo) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type AddItemInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    uint32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Items []uint32 `protobuf:"varint,2,rep,packed,name=items,proto3" json:"items,omitempty"`
}

func (x *AddItemInfo) Reset() {
	*x = AddItemInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddItemInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddItemInfo) ProtoMessage() {}

func (x *AddItemInfo) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddItemInfo.ProtoReflect.Descriptor instead.
func (*AddItemInfo) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{5}
}

func (x *AddItemInfo) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AddItemInfo) GetItems() []uint32 {
	if x != nil {
		return x.Items
	}
	return nil
}

type RoomIDInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoomID uint32 `protobuf:"varint,1,opt,name=RoomID,proto3" json:"RoomID,omitempty"`
}

func (x *RoomIDInfo) Reset() {
	*x = RoomIDInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoomIDInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoomIDInfo) ProtoMessage() {}

func (x *RoomIDInfo) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoomIDInfo.ProtoReflect.Descriptor instead.
func (*RoomIDInfo) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{6}
}

func (x *RoomIDInfo) GetRoomID() uint32 {
	if x != nil {
		return x.RoomID
	}
	return 0
}

type RoomInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoomID        uint32 `protobuf:"varint,1,opt,name=RoomID,proto3" json:"RoomID,omitempty"`
	MaxUserNumber uint32 `protobuf:"varint,2,opt,name=MaxUserNumber,proto3" json:"MaxUserNumber,omitempty"`
	GameCount     uint32 `protobuf:"varint,3,opt,name=GameCount,proto3" json:"GameCount,omitempty"`
	UserNumber    uint32 `protobuf:"varint,4,opt,name=UserNumber,proto3" json:"UserNumber,omitempty"`
	RoomOwner     uint32 `protobuf:"varint,5,opt,name=RoomOwner,proto3" json:"RoomOwner,omitempty"`
	RoomWait      bool   `protobuf:"varint,6,opt,name=RoomWait,proto3" json:"RoomWait,omitempty"`
}

func (x *RoomInfo) Reset() {
	*x = RoomInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoomInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoomInfo) ProtoMessage() {}

func (x *RoomInfo) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoomInfo.ProtoReflect.Descriptor instead.
func (*RoomInfo) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{7}
}

func (x *RoomInfo) GetRoomID() uint32 {
	if x != nil {
		return x.RoomID
	}
	return 0
}

func (x *RoomInfo) GetMaxUserNumber() uint32 {
	if x != nil {
		return x.MaxUserNumber
	}
	return 0
}

func (x *RoomInfo) GetGameCount() uint32 {
	if x != nil {
		return x.GameCount
	}
	return 0
}

func (x *RoomInfo) GetUserNumber() uint32 {
	if x != nil {
		return x.UserNumber
	}
	return 0
}

func (x *RoomInfo) GetRoomOwner() uint32 {
	if x != nil {
		return x.RoomOwner
	}
	return 0
}

func (x *RoomInfo) GetRoomWait() bool {
	if x != nil {
		return x.RoomWait
	}
	return false
}

type AllRoomInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AllRoomInfo []*RoomInfo `protobuf:"bytes,1,rep,name=AllRoomInfo,proto3" json:"AllRoomInfo,omitempty"`
}

func (x *AllRoomInfo) Reset() {
	*x = AllRoomInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllRoomInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllRoomInfo) ProtoMessage() {}

func (x *AllRoomInfo) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllRoomInfo.ProtoReflect.Descriptor instead.
func (*AllRoomInfo) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{8}
}

func (x *AllRoomInfo) GetAllRoomInfo() []*RoomInfo {
	if x != nil {
		return x.AllRoomInfo
	}
	return nil
}

var File_game_proto protoreflect.FileDescriptor

var file_game_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7b, 0x0a, 0x0d, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x6f,
	0x6c, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x67, 0x6f, 0x6c, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x64, 0x69, 0x61, 0x6d, 0x6f, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x07, 0x64, 0x69, 0x61, 0x6d, 0x6f, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x70, 0x70, 0x6c,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x61, 0x70, 0x70, 0x6c, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x62, 0x61, 0x6e, 0x61, 0x6e, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06,
	0x62, 0x61, 0x6e, 0x61, 0x6e, 0x61, 0x22, 0x1c, 0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x6b, 0x0a, 0x15, 0x55, 0x73, 0x65, 0x72, 0x49, 0x74, 0x65, 0x6d,
	0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x67, 0x6f, 0x6c, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x67, 0x6f, 0x6c,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x69, 0x61, 0x6d, 0x6f, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x64, 0x69, 0x61, 0x6d, 0x6f, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x69,
	0x74, 0x65, 0x6d, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x22, 0x36, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x47, 0x6f, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x36, 0x0a, 0x0e, 0x41, 0x64, 0x64,
	0x44, 0x69, 0x61, 0x6d, 0x6f, 0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x22, 0x33, 0x0a, 0x0b, 0x41, 0x64, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0d, 0x52,
	0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x24, 0x0a, 0x0a, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x44,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x44, 0x22, 0xc0, 0x01, 0x0a,
	0x08, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x6f, 0x6f,
	0x6d, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x52, 0x6f, 0x6f, 0x6d, 0x49,
	0x44, 0x12, 0x24, 0x0a, 0x0d, 0x4d, 0x61, 0x78, 0x55, 0x73, 0x65, 0x72, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x4d, 0x61, 0x78, 0x55, 0x73, 0x65,
	0x72, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x47, 0x61, 0x6d, 0x65, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x47, 0x61, 0x6d, 0x65,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x52, 0x6f, 0x6f, 0x6d, 0x4f, 0x77, 0x6e,
	0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x52, 0x6f, 0x6f, 0x6d, 0x4f, 0x77,
	0x6e, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x52, 0x6f, 0x6f, 0x6d, 0x57, 0x61, 0x69, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x52, 0x6f, 0x6f, 0x6d, 0x57, 0x61, 0x69, 0x74, 0x22,
	0x3a, 0x0a, 0x0b, 0x41, 0x6c, 0x6c, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2b,
	0x0a, 0x0b, 0x41, 0x6c, 0x6c, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0b,
	0x41, 0x6c, 0x6c, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x2a, 0x1d, 0x0a, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x10, 0x00, 0x12, 0x0a,
	0x0a, 0x06, 0x42, 0x61, 0x6e, 0x61, 0x6e, 0x61, 0x10, 0x01, 0x32, 0x88, 0x04, 0x0a, 0x04, 0x47,
	0x61, 0x6d, 0x65, 0x12, 0x39, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x0e, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x74, 0x65,
	0x6d, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x16, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x74, 0x65,
	0x6d, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x0b, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x49, 0x6e, 0x66, 0x6f, 0x1a,
	0x16, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x47, 0x6f,
	0x6c, 0x64, 0x12, 0x0f, 0x2e, 0x41, 0x64, 0x64, 0x47, 0x6f, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x35, 0x0a, 0x0a, 0x41,
	0x64, 0x64, 0x44, 0x69, 0x61, 0x6d, 0x6f, 0x6e, 0x64, 0x12, 0x0f, 0x2e, 0x41, 0x64, 0x64, 0x44,
	0x69, 0x61, 0x6d, 0x6f, 0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x2f, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x0c, 0x2e,
	0x41, 0x64, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x12, 0x35, 0x0a, 0x0d, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x41, 0x6c, 0x6c,
	0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e, 0x41,
	0x6c, 0x6c, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2f, 0x0a, 0x0a, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x09, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x49,
	0x6e, 0x66, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x24, 0x0a, 0x0a, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x0b, 0x2e, 0x52, 0x6f, 0x6f, 0x6d,
	0x49, 0x44, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x09, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x2f, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x6d, 0x12,
	0x09, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x31, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x6d,
	0x12, 0x0b, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x44, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_game_proto_rawDescOnce sync.Once
	file_game_proto_rawDescData = file_game_proto_rawDesc
)

func file_game_proto_rawDescGZIP() []byte {
	file_game_proto_rawDescOnce.Do(func() {
		file_game_proto_rawDescData = protoimpl.X.CompressGZIP(file_game_proto_rawDescData)
	})
	return file_game_proto_rawDescData
}

var file_game_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_game_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_game_proto_goTypes = []interface{}{
	(Type)(0),                     // 0: Type
	(*UserItemsInfo)(nil),         // 1: UserItemsInfo
	(*UserIDInfo)(nil),            // 2: UserIDInfo
	(*UserItemsInfoResponse)(nil), // 3: UserItemsInfoResponse
	(*AddGoldRequest)(nil),        // 4: AddGoldRequest
	(*AddDiamondInfo)(nil),        // 5: AddDiamondInfo
	(*AddItemInfo)(nil),           // 6: AddItemInfo
	(*RoomIDInfo)(nil),            // 7: RoomIDInfo
	(*RoomInfo)(nil),              // 8: RoomInfo
	(*AllRoomInfo)(nil),           // 9: AllRoomInfo
	(*emptypb.Empty)(nil),         // 10: google.protobuf.Empty
}
var file_game_proto_depIdxs = []int32{
	8,  // 0: AllRoomInfo.AllRoomInfo:type_name -> RoomInfo
	1,  // 1: Game.CreateUserItems:input_type -> UserItemsInfo
	2,  // 2: Game.GetUserItemsInfo:input_type -> UserIDInfo
	4,  // 3: Game.AddGold:input_type -> AddGoldRequest
	5,  // 4: Game.AddDiamond:input_type -> AddDiamondInfo
	6,  // 5: Game.AddItem:input_type -> AddItemInfo
	10, // 6: Game.SearchAllRoom:input_type -> google.protobuf.Empty
	8,  // 7: Game.CreateRoom:input_type -> RoomInfo
	7,  // 8: Game.SearchRoom:input_type -> RoomIDInfo
	8,  // 9: Game.UpdateRoom:input_type -> RoomInfo
	7,  // 10: Game.DeleteRoom:input_type -> RoomIDInfo
	3,  // 11: Game.CreateUserItems:output_type -> UserItemsInfoResponse
	3,  // 12: Game.GetUserItemsInfo:output_type -> UserItemsInfoResponse
	10, // 13: Game.AddGold:output_type -> google.protobuf.Empty
	10, // 14: Game.AddDiamond:output_type -> google.protobuf.Empty
	10, // 15: Game.AddItem:output_type -> google.protobuf.Empty
	9,  // 16: Game.SearchAllRoom:output_type -> AllRoomInfo
	10, // 17: Game.CreateRoom:output_type -> google.protobuf.Empty
	8,  // 18: Game.SearchRoom:output_type -> RoomInfo
	10, // 19: Game.UpdateRoom:output_type -> google.protobuf.Empty
	10, // 20: Game.DeleteRoom:output_type -> google.protobuf.Empty
	11, // [11:21] is the sub-list for method output_type
	1,  // [1:11] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_game_proto_init() }
func file_game_proto_init() {
	if File_game_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_game_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserItemsInfo); i {
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
		file_game_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserIDInfo); i {
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
		file_game_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserItemsInfoResponse); i {
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
		file_game_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddGoldRequest); i {
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
		file_game_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddDiamondInfo); i {
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
		file_game_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddItemInfo); i {
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
		file_game_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoomIDInfo); i {
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
		file_game_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoomInfo); i {
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
		file_game_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllRoomInfo); i {
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
			RawDescriptor: file_game_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_game_proto_goTypes,
		DependencyIndexes: file_game_proto_depIdxs,
		EnumInfos:         file_game_proto_enumTypes,
		MessageInfos:      file_game_proto_msgTypes,
	}.Build()
	File_game_proto = out.File
	file_game_proto_rawDesc = nil
	file_game_proto_goTypes = nil
	file_game_proto_depIdxs = nil
}
