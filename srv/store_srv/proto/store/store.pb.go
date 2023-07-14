// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: store.proto

package store

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

type SelectGoodsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page     uint32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize uint32 `protobuf:"varint,2,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
}

func (x *SelectGoodsReq) Reset() {
	*x = SelectGoodsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SelectGoodsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SelectGoodsReq) ProtoMessage() {}

func (x *SelectGoodsReq) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SelectGoodsReq.ProtoReflect.Descriptor instead.
func (*SelectGoodsReq) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{0}
}

func (x *SelectGoodsReq) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SelectGoodsReq) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type SelectGoodsRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnKnow []uint32 `protobuf:"varint,1,rep,packed,name=unKnow,proto3" json:"unKnow,omitempty"`
}

func (x *SelectGoodsRsp) Reset() {
	*x = SelectGoodsRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SelectGoodsRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SelectGoodsRsp) ProtoMessage() {}

func (x *SelectGoodsRsp) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SelectGoodsRsp.ProtoReflect.Descriptor instead.
func (*SelectGoodsRsp) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{1}
}

func (x *SelectGoodsRsp) GetUnKnow() []uint32 {
	if x != nil {
		return x.UnKnow
	}
	return nil
}

type BuyGoodReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` //购买者
	GoodId uint32 `protobuf:"varint,2,opt,name=good_id,json=goodId,proto3" json:"good_id,omitempty"` //购买商品id
	Count  uint32 `protobuf:"varint,3,opt,name=count,proto3" json:"count,omitempty"`                 //购买数量
}

func (x *BuyGoodReq) Reset() {
	*x = BuyGoodReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuyGoodReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuyGoodReq) ProtoMessage() {}

func (x *BuyGoodReq) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuyGoodReq.ProtoReflect.Descriptor instead.
func (*BuyGoodReq) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{2}
}

func (x *BuyGoodReq) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *BuyGoodReq) GetGoodId() uint32 {
	if x != nil {
		return x.GoodId
	}
	return 0
}

func (x *BuyGoodReq) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type SelectTradeItemsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page     uint32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize uint32 `protobuf:"varint,2,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
}

func (x *SelectTradeItemsReq) Reset() {
	*x = SelectTradeItemsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SelectTradeItemsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SelectTradeItemsReq) ProtoMessage() {}

func (x *SelectTradeItemsReq) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SelectTradeItemsReq.ProtoReflect.Descriptor instead.
func (*SelectTradeItemsReq) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{3}
}

func (x *SelectTradeItemsReq) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SelectTradeItemsReq) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type SelectTradeItemsRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnKnow []uint32 `protobuf:"varint,1,rep,packed,name=unKnow,proto3" json:"unKnow,omitempty"`
}

func (x *SelectTradeItemsRsp) Reset() {
	*x = SelectTradeItemsRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SelectTradeItemsRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SelectTradeItemsRsp) ProtoMessage() {}

func (x *SelectTradeItemsRsp) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SelectTradeItemsRsp.ProtoReflect.Descriptor instead.
func (*SelectTradeItemsRsp) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{4}
}

func (x *SelectTradeItemsRsp) GetUnKnow() []uint32 {
	if x != nil {
		return x.UnKnow
	}
	return nil
}

type PushTradeItemReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` //售卖者
	ItemId uint32 `protobuf:"varint,2,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"` //物品
	Price  uint32 `protobuf:"varint,3,opt,name=price,proto3" json:"price,omitempty"`                 //售卖价格
	Count  uint32 `protobuf:"varint,4,opt,name=count,proto3" json:"count,omitempty"`                 //售卖数量
}

func (x *PushTradeItemReq) Reset() {
	*x = PushTradeItemReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushTradeItemReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushTradeItemReq) ProtoMessage() {}

func (x *PushTradeItemReq) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushTradeItemReq.ProtoReflect.Descriptor instead.
func (*PushTradeItemReq) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{5}
}

func (x *PushTradeItemReq) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *PushTradeItemReq) GetItemId() uint32 {
	if x != nil {
		return x.ItemId
	}
	return 0
}

func (x *PushTradeItemReq) GetPrice() uint32 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *PushTradeItemReq) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type DownTradeItemReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TradeId uint32 `protobuf:"varint,1,opt,name=trade_id,json=tradeId,proto3" json:"trade_id,omitempty"` //售卖唯一ID
}

func (x *DownTradeItemReq) Reset() {
	*x = DownTradeItemReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownTradeItemReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownTradeItemReq) ProtoMessage() {}

func (x *DownTradeItemReq) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownTradeItemReq.ProtoReflect.Descriptor instead.
func (*DownTradeItemReq) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{6}
}

func (x *DownTradeItemReq) GetTradeId() uint32 {
	if x != nil {
		return x.TradeId
	}
	return 0
}

type BuyTradeItemReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TradeId uint32 `protobuf:"varint,1,opt,name=trade_id,json=tradeId,proto3" json:"trade_id,omitempty"` //售卖ID
}

func (x *BuyTradeItemReq) Reset() {
	*x = BuyTradeItemReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuyTradeItemReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuyTradeItemReq) ProtoMessage() {}

func (x *BuyTradeItemReq) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuyTradeItemReq.ProtoReflect.Descriptor instead.
func (*BuyTradeItemReq) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{7}
}

func (x *BuyTradeItemReq) GetTradeId() uint32 {
	if x != nil {
		return x.TradeId
	}
	return 0
}

var File_store_proto protoreflect.FileDescriptor

var file_store_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x40, 0x0a, 0x0e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x47, 0x6f, 0x6f, 0x64, 0x73,
	0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53,
	0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53,
	0x69, 0x7a, 0x65, 0x22, 0x28, 0x0a, 0x0e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x47, 0x6f, 0x6f,
	0x64, 0x73, 0x52, 0x73, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x6e, 0x4b, 0x6e, 0x6f, 0x77, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x6e, 0x4b, 0x6e, 0x6f, 0x77, 0x22, 0x54, 0x0a,
	0x0a, 0x42, 0x75, 0x79, 0x47, 0x6f, 0x6f, 0x64, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x67, 0x6f, 0x6f, 0x64, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x67, 0x6f, 0x6f, 0x64, 0x49, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x22, 0x45, 0x0a, 0x13, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x54, 0x72, 0x61,
	0x64, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x2d, 0x0a, 0x13, 0x53, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x73,
	0x70, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x6e, 0x4b, 0x6e, 0x6f, 0x77, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0d, 0x52, 0x06, 0x75, 0x6e, 0x4b, 0x6e, 0x6f, 0x77, 0x22, 0x70, 0x0a, 0x10, 0x50, 0x75, 0x73,
	0x68, 0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x2d, 0x0a, 0x10, 0x44,
	0x6f, 0x77, 0x6e, 0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x12,
	0x19, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x07, 0x74, 0x72, 0x61, 0x64, 0x65, 0x49, 0x64, 0x22, 0x2c, 0x0a, 0x0f, 0x42, 0x75,
	0x79, 0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x19, 0x0a,
	0x08, 0x74, 0x72, 0x61, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x07, 0x74, 0x72, 0x61, 0x64, 0x65, 0x49, 0x64, 0x32, 0x8b, 0x03, 0x0a, 0x05, 0x53, 0x74, 0x6f,
	0x72, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x47, 0x6f, 0x6f, 0x64,
	0x73, 0x12, 0x15, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x47, 0x6f, 0x6f, 0x64, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x15, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x52, 0x73, 0x70, 0x12,
	0x35, 0x0a, 0x08, 0x42, 0x75, 0x79, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x12, 0x11, 0x2e, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x42, 0x75, 0x79, 0x47, 0x6f, 0x6f, 0x64, 0x52, 0x65, 0x71, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x4a, 0x0a, 0x10, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x1a, 0x2e, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x74,
	0x65, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x1a, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x53,
	0x65, 0x6c, 0x65, 0x63, 0x74, 0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52,
	0x73, 0x70, 0x12, 0x40, 0x0a, 0x0d, 0x50, 0x75, 0x73, 0x68, 0x54, 0x72, 0x61, 0x64, 0x65, 0x49,
	0x74, 0x65, 0x6d, 0x12, 0x17, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x50, 0x75, 0x73, 0x68,
	0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x12, 0x40, 0x0a, 0x0d, 0x44, 0x6f, 0x77, 0x6e, 0x54, 0x72, 0x61, 0x64,
	0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x17, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x44, 0x6f,
	0x77, 0x6e, 0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x3e, 0x0a, 0x0c, 0x42, 0x75, 0x79, 0x54, 0x72, 0x61,
	0x64, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x16, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x42,
	0x75, 0x79, 0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x3b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_store_proto_rawDescOnce sync.Once
	file_store_proto_rawDescData = file_store_proto_rawDesc
)

func file_store_proto_rawDescGZIP() []byte {
	file_store_proto_rawDescOnce.Do(func() {
		file_store_proto_rawDescData = protoimpl.X.CompressGZIP(file_store_proto_rawDescData)
	})
	return file_store_proto_rawDescData
}

var file_store_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_store_proto_goTypes = []interface{}{
	(*SelectGoodsReq)(nil),      // 0: store.SelectGoodsReq
	(*SelectGoodsRsp)(nil),      // 1: store.SelectGoodsRsp
	(*BuyGoodReq)(nil),          // 2: store.BuyGoodReq
	(*SelectTradeItemsReq)(nil), // 3: store.SelectTradeItemsReq
	(*SelectTradeItemsRsp)(nil), // 4: store.SelectTradeItemsRsp
	(*PushTradeItemReq)(nil),    // 5: store.PushTradeItemReq
	(*DownTradeItemReq)(nil),    // 6: store.DownTradeItemReq
	(*BuyTradeItemReq)(nil),     // 7: store.BuyTradeItemReq
	(*emptypb.Empty)(nil),       // 8: google.protobuf.Empty
}
var file_store_proto_depIdxs = []int32{
	0, // 0: store.Store.SelectGoods:input_type -> store.SelectGoodsReq
	2, // 1: store.Store.BuyGoods:input_type -> store.BuyGoodReq
	3, // 2: store.Store.SelectTradeItems:input_type -> store.SelectTradeItemsReq
	5, // 3: store.Store.PushTradeItem:input_type -> store.PushTradeItemReq
	6, // 4: store.Store.DownTradeItem:input_type -> store.DownTradeItemReq
	7, // 5: store.Store.BuyTradeItem:input_type -> store.BuyTradeItemReq
	1, // 6: store.Store.SelectGoods:output_type -> store.SelectGoodsRsp
	8, // 7: store.Store.BuyGoods:output_type -> google.protobuf.Empty
	4, // 8: store.Store.SelectTradeItems:output_type -> store.SelectTradeItemsRsp
	8, // 9: store.Store.PushTradeItem:output_type -> google.protobuf.Empty
	8, // 10: store.Store.DownTradeItem:output_type -> google.protobuf.Empty
	8, // 11: store.Store.BuyTradeItem:output_type -> google.protobuf.Empty
	6, // [6:12] is the sub-list for method output_type
	0, // [0:6] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_store_proto_init() }
func file_store_proto_init() {
	if File_store_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_store_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SelectGoodsReq); i {
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
		file_store_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SelectGoodsRsp); i {
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
		file_store_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuyGoodReq); i {
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
		file_store_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SelectTradeItemsReq); i {
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
		file_store_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SelectTradeItemsRsp); i {
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
		file_store_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushTradeItemReq); i {
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
		file_store_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownTradeItemReq); i {
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
		file_store_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuyTradeItemReq); i {
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
			RawDescriptor: file_store_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_store_proto_goTypes,
		DependencyIndexes: file_store_proto_depIdxs,
		MessageInfos:      file_store_proto_msgTypes,
	}.Build()
	File_store_proto = out.File
	file_store_proto_rawDesc = nil
	file_store_proto_goTypes = nil
	file_store_proto_depIdxs = nil
}
