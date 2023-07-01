// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package game

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GameClient is the client API for Game service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GameClient interface {
	// 获得连接的服务器信息
	GetConnData(ctx context.Context, in *UserIDInfo, opts ...grpc.CallOption) (*ConnResponse, error)
	// 记录连接的服务器信息
	RecordConnData(ctx context.Context, in *RecordConnInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 获得排行榜信息
	GetRanks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*RanksResponse, error)
	// 更新排行榜
	UpdateRanks(ctx context.Context, in *UpdateRanksInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 新建用户道具表
	CreateUserItems(ctx context.Context, in *UserItemsInfo, opts ...grpc.CallOption) (*UserItemsInfoResponse, error)
	// 查询用户的金币钻石及道具
	GetUserItemsInfo(ctx context.Context, in *UserIDInfo, opts ...grpc.CallOption) (*UserItemsInfoResponse, error)
	// 增加金币
	AddGold(ctx context.Context, in *AddGoldRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 增加钻石
	AddDiamond(ctx context.Context, in *AddDiamondInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 增加道具(道具类型应该区别)
	AddItem(ctx context.Context, in *AddItemInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 使用金币
	UseGold(ctx context.Context, in *UseGoldRequest, opts ...grpc.CallOption) (*IsOK, error)
	// 使用钻石
	UseDiamond(ctx context.Context, in *UseDiamondInfo, opts ...grpc.CallOption) (*IsOK, error)
	// 使用道具(道具类型应该区别)
	UseItem(ctx context.Context, in *UseItemInfo, opts ...grpc.CallOption) (*IsOK, error)
	// 查询所有房间
	SearchAllRoom(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AllRoomInfo, error)
	// 创建房间
	CreateRoom(ctx context.Context, in *RoomInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 进入房间
	UserIntoRoom(ctx context.Context, in *UserIntoRoomInfo, opts ...grpc.CallOption) (*IntoRoomRsp, error)
	// 查询房间
	SearchRoom(ctx context.Context, in *RoomIDInfo, opts ...grpc.CallOption) (*RoomInfo, error)
	// 删除房间
	DeleteRoom(ctx context.Context, in *RoomIDInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 退出房间
	QuitRoom(ctx context.Context, in *QuitRoomInfo, opts ...grpc.CallOption) (*QuitRsp, error)
	// 房主更新房间信息
	UpdateRoom(ctx context.Context, in *UpdateRoomInfo, opts ...grpc.CallOption) (*RoomInfo, error)
	// 更新用户准备状态
	UpdateUserReadyState(ctx context.Context, in *ReadyStateInfo, opts ...grpc.CallOption) (*RoomInfo, error)
	// 房主开始游戏
	BeginGame(ctx context.Context, in *BeginGameInfo, opts ...grpc.CallOption) (*BeginGameRsp, error)
	// 回到房间
	BackRoom(ctx context.Context, in *RoomIDInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type gameClient struct {
	cc grpc.ClientConnInterface
}

func NewGameClient(cc grpc.ClientConnInterface) GameClient {
	return &gameClient{cc}
}

func (c *gameClient) GetConnData(ctx context.Context, in *UserIDInfo, opts ...grpc.CallOption) (*ConnResponse, error) {
	out := new(ConnResponse)
	err := c.cc.Invoke(ctx, "/game.Game/GetConnData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) RecordConnData(ctx context.Context, in *RecordConnInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/game.Game/RecordConnData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) GetRanks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*RanksResponse, error) {
	out := new(RanksResponse)
	err := c.cc.Invoke(ctx, "/game.Game/GetRanks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) UpdateRanks(ctx context.Context, in *UpdateRanksInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/game.Game/UpdateRanks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) CreateUserItems(ctx context.Context, in *UserItemsInfo, opts ...grpc.CallOption) (*UserItemsInfoResponse, error) {
	out := new(UserItemsInfoResponse)
	err := c.cc.Invoke(ctx, "/game.Game/CreateUserItems", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) GetUserItemsInfo(ctx context.Context, in *UserIDInfo, opts ...grpc.CallOption) (*UserItemsInfoResponse, error) {
	out := new(UserItemsInfoResponse)
	err := c.cc.Invoke(ctx, "/game.Game/GetUserItemsInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) AddGold(ctx context.Context, in *AddGoldRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/game.Game/AddGold", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) AddDiamond(ctx context.Context, in *AddDiamondInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/game.Game/AddDiamond", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) AddItem(ctx context.Context, in *AddItemInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/game.Game/AddItem", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) UseGold(ctx context.Context, in *UseGoldRequest, opts ...grpc.CallOption) (*IsOK, error) {
	out := new(IsOK)
	err := c.cc.Invoke(ctx, "/game.Game/UseGold", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) UseDiamond(ctx context.Context, in *UseDiamondInfo, opts ...grpc.CallOption) (*IsOK, error) {
	out := new(IsOK)
	err := c.cc.Invoke(ctx, "/game.Game/UseDiamond", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) UseItem(ctx context.Context, in *UseItemInfo, opts ...grpc.CallOption) (*IsOK, error) {
	out := new(IsOK)
	err := c.cc.Invoke(ctx, "/game.Game/UseItem", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) SearchAllRoom(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AllRoomInfo, error) {
	out := new(AllRoomInfo)
	err := c.cc.Invoke(ctx, "/game.Game/SearchAllRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) CreateRoom(ctx context.Context, in *RoomInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/game.Game/CreateRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) UserIntoRoom(ctx context.Context, in *UserIntoRoomInfo, opts ...grpc.CallOption) (*IntoRoomRsp, error) {
	out := new(IntoRoomRsp)
	err := c.cc.Invoke(ctx, "/game.Game/UserIntoRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) SearchRoom(ctx context.Context, in *RoomIDInfo, opts ...grpc.CallOption) (*RoomInfo, error) {
	out := new(RoomInfo)
	err := c.cc.Invoke(ctx, "/game.Game/SearchRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) DeleteRoom(ctx context.Context, in *RoomIDInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/game.Game/DeleteRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) QuitRoom(ctx context.Context, in *QuitRoomInfo, opts ...grpc.CallOption) (*QuitRsp, error) {
	out := new(QuitRsp)
	err := c.cc.Invoke(ctx, "/game.Game/QuitRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) UpdateRoom(ctx context.Context, in *UpdateRoomInfo, opts ...grpc.CallOption) (*RoomInfo, error) {
	out := new(RoomInfo)
	err := c.cc.Invoke(ctx, "/game.Game/UpdateRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) UpdateUserReadyState(ctx context.Context, in *ReadyStateInfo, opts ...grpc.CallOption) (*RoomInfo, error) {
	out := new(RoomInfo)
	err := c.cc.Invoke(ctx, "/game.Game/UpdateUserReadyState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) BeginGame(ctx context.Context, in *BeginGameInfo, opts ...grpc.CallOption) (*BeginGameRsp, error) {
	out := new(BeginGameRsp)
	err := c.cc.Invoke(ctx, "/game.Game/BeginGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) BackRoom(ctx context.Context, in *RoomIDInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/game.Game/BackRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GameServer is the server API for Game service.
// All implementations must embed UnimplementedGameServer
// for forward compatibility
type GameServer interface {
	// 获得连接的服务器信息
	GetConnData(context.Context, *UserIDInfo) (*ConnResponse, error)
	// 记录连接的服务器信息
	RecordConnData(context.Context, *RecordConnInfo) (*emptypb.Empty, error)
	// 获得排行榜信息
	GetRanks(context.Context, *emptypb.Empty) (*RanksResponse, error)
	// 更新排行榜
	UpdateRanks(context.Context, *UpdateRanksInfo) (*emptypb.Empty, error)
	// 新建用户道具表
	CreateUserItems(context.Context, *UserItemsInfo) (*UserItemsInfoResponse, error)
	// 查询用户的金币钻石及道具
	GetUserItemsInfo(context.Context, *UserIDInfo) (*UserItemsInfoResponse, error)
	// 增加金币
	AddGold(context.Context, *AddGoldRequest) (*emptypb.Empty, error)
	// 增加钻石
	AddDiamond(context.Context, *AddDiamondInfo) (*emptypb.Empty, error)
	// 增加道具(道具类型应该区别)
	AddItem(context.Context, *AddItemInfo) (*emptypb.Empty, error)
	// 使用金币
	UseGold(context.Context, *UseGoldRequest) (*IsOK, error)
	// 使用钻石
	UseDiamond(context.Context, *UseDiamondInfo) (*IsOK, error)
	// 使用道具(道具类型应该区别)
	UseItem(context.Context, *UseItemInfo) (*IsOK, error)
	// 查询所有房间
	SearchAllRoom(context.Context, *emptypb.Empty) (*AllRoomInfo, error)
	// 创建房间
	CreateRoom(context.Context, *RoomInfo) (*emptypb.Empty, error)
	// 进入房间
	UserIntoRoom(context.Context, *UserIntoRoomInfo) (*IntoRoomRsp, error)
	// 查询房间
	SearchRoom(context.Context, *RoomIDInfo) (*RoomInfo, error)
	// 删除房间
	DeleteRoom(context.Context, *RoomIDInfo) (*emptypb.Empty, error)
	// 退出房间
	QuitRoom(context.Context, *QuitRoomInfo) (*QuitRsp, error)
	// 房主更新房间信息
	UpdateRoom(context.Context, *UpdateRoomInfo) (*RoomInfo, error)
	// 更新用户准备状态
	UpdateUserReadyState(context.Context, *ReadyStateInfo) (*RoomInfo, error)
	// 房主开始游戏
	BeginGame(context.Context, *BeginGameInfo) (*BeginGameRsp, error)
	// 回到房间
	BackRoom(context.Context, *RoomIDInfo) (*emptypb.Empty, error)
	mustEmbedUnimplementedGameServer()
}

// UnimplementedGameServer must be embedded to have forward compatible implementations.
type UnimplementedGameServer struct {
}

func (UnimplementedGameServer) GetConnData(context.Context, *UserIDInfo) (*ConnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConnData not implemented")
}
func (UnimplementedGameServer) RecordConnData(context.Context, *RecordConnInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordConnData not implemented")
}
func (UnimplementedGameServer) GetRanks(context.Context, *emptypb.Empty) (*RanksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRanks not implemented")
}
func (UnimplementedGameServer) UpdateRanks(context.Context, *UpdateRanksInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRanks not implemented")
}
func (UnimplementedGameServer) CreateUserItems(context.Context, *UserItemsInfo) (*UserItemsInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserItems not implemented")
}
func (UnimplementedGameServer) GetUserItemsInfo(context.Context, *UserIDInfo) (*UserItemsInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserItemsInfo not implemented")
}
func (UnimplementedGameServer) AddGold(context.Context, *AddGoldRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddGold not implemented")
}
func (UnimplementedGameServer) AddDiamond(context.Context, *AddDiamondInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddDiamond not implemented")
}
func (UnimplementedGameServer) AddItem(context.Context, *AddItemInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddItem not implemented")
}
func (UnimplementedGameServer) UseGold(context.Context, *UseGoldRequest) (*IsOK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseGold not implemented")
}
func (UnimplementedGameServer) UseDiamond(context.Context, *UseDiamondInfo) (*IsOK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseDiamond not implemented")
}
func (UnimplementedGameServer) UseItem(context.Context, *UseItemInfo) (*IsOK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseItem not implemented")
}
func (UnimplementedGameServer) SearchAllRoom(context.Context, *emptypb.Empty) (*AllRoomInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchAllRoom not implemented")
}
func (UnimplementedGameServer) CreateRoom(context.Context, *RoomInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoom not implemented")
}
func (UnimplementedGameServer) UserIntoRoom(context.Context, *UserIntoRoomInfo) (*IntoRoomRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserIntoRoom not implemented")
}
func (UnimplementedGameServer) SearchRoom(context.Context, *RoomIDInfo) (*RoomInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchRoom not implemented")
}
func (UnimplementedGameServer) DeleteRoom(context.Context, *RoomIDInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRoom not implemented")
}
func (UnimplementedGameServer) QuitRoom(context.Context, *QuitRoomInfo) (*QuitRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QuitRoom not implemented")
}
func (UnimplementedGameServer) UpdateRoom(context.Context, *UpdateRoomInfo) (*RoomInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRoom not implemented")
}
func (UnimplementedGameServer) UpdateUserReadyState(context.Context, *ReadyStateInfo) (*RoomInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserReadyState not implemented")
}
func (UnimplementedGameServer) BeginGame(context.Context, *BeginGameInfo) (*BeginGameRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BeginGame not implemented")
}
func (UnimplementedGameServer) BackRoom(context.Context, *RoomIDInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BackRoom not implemented")
}
func (UnimplementedGameServer) mustEmbedUnimplementedGameServer() {}

// UnsafeGameServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GameServer will
// result in compilation errors.
type UnsafeGameServer interface {
	mustEmbedUnimplementedGameServer()
}

func RegisterGameServer(s grpc.ServiceRegistrar, srv GameServer) {
	s.RegisterService(&Game_ServiceDesc, srv)
}

func _Game_GetConnData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIDInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).GetConnData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/GetConnData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).GetConnData(ctx, req.(*UserIDInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_RecordConnData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordConnInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).RecordConnData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/RecordConnData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).RecordConnData(ctx, req.(*RecordConnInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_GetRanks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).GetRanks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/GetRanks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).GetRanks(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_UpdateRanks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRanksInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).UpdateRanks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/UpdateRanks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).UpdateRanks(ctx, req.(*UpdateRanksInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_CreateUserItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserItemsInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).CreateUserItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/CreateUserItems",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).CreateUserItems(ctx, req.(*UserItemsInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_GetUserItemsInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIDInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).GetUserItemsInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/GetUserItemsInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).GetUserItemsInfo(ctx, req.(*UserIDInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_AddGold_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddGoldRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).AddGold(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/AddGold",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).AddGold(ctx, req.(*AddGoldRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_AddDiamond_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddDiamondInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).AddDiamond(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/AddDiamond",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).AddDiamond(ctx, req.(*AddDiamondInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_AddItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddItemInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).AddItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/AddItem",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).AddItem(ctx, req.(*AddItemInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_UseGold_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UseGoldRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).UseGold(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/UseGold",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).UseGold(ctx, req.(*UseGoldRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_UseDiamond_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UseDiamondInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).UseDiamond(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/UseDiamond",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).UseDiamond(ctx, req.(*UseDiamondInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_UseItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UseItemInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).UseItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/UseItem",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).UseItem(ctx, req.(*UseItemInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_SearchAllRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).SearchAllRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/SearchAllRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).SearchAllRoom(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_CreateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).CreateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/CreateRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).CreateRoom(ctx, req.(*RoomInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_UserIntoRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIntoRoomInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).UserIntoRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/UserIntoRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).UserIntoRoom(ctx, req.(*UserIntoRoomInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_SearchRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomIDInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).SearchRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/SearchRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).SearchRoom(ctx, req.(*RoomIDInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_DeleteRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomIDInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).DeleteRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/DeleteRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).DeleteRoom(ctx, req.(*RoomIDInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_QuitRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuitRoomInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).QuitRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/QuitRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).QuitRoom(ctx, req.(*QuitRoomInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_UpdateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRoomInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).UpdateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/UpdateRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).UpdateRoom(ctx, req.(*UpdateRoomInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_UpdateUserReadyState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadyStateInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).UpdateUserReadyState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/UpdateUserReadyState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).UpdateUserReadyState(ctx, req.(*ReadyStateInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_BeginGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BeginGameInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).BeginGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/BeginGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).BeginGame(ctx, req.(*BeginGameInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_BackRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomIDInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).BackRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.Game/BackRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).BackRoom(ctx, req.(*RoomIDInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// Game_ServiceDesc is the grpc.ServiceDesc for Game service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Game_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "game.Game",
	HandlerType: (*GameServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetConnData",
			Handler:    _Game_GetConnData_Handler,
		},
		{
			MethodName: "RecordConnData",
			Handler:    _Game_RecordConnData_Handler,
		},
		{
			MethodName: "GetRanks",
			Handler:    _Game_GetRanks_Handler,
		},
		{
			MethodName: "UpdateRanks",
			Handler:    _Game_UpdateRanks_Handler,
		},
		{
			MethodName: "CreateUserItems",
			Handler:    _Game_CreateUserItems_Handler,
		},
		{
			MethodName: "GetUserItemsInfo",
			Handler:    _Game_GetUserItemsInfo_Handler,
		},
		{
			MethodName: "AddGold",
			Handler:    _Game_AddGold_Handler,
		},
		{
			MethodName: "AddDiamond",
			Handler:    _Game_AddDiamond_Handler,
		},
		{
			MethodName: "AddItem",
			Handler:    _Game_AddItem_Handler,
		},
		{
			MethodName: "UseGold",
			Handler:    _Game_UseGold_Handler,
		},
		{
			MethodName: "UseDiamond",
			Handler:    _Game_UseDiamond_Handler,
		},
		{
			MethodName: "UseItem",
			Handler:    _Game_UseItem_Handler,
		},
		{
			MethodName: "SearchAllRoom",
			Handler:    _Game_SearchAllRoom_Handler,
		},
		{
			MethodName: "CreateRoom",
			Handler:    _Game_CreateRoom_Handler,
		},
		{
			MethodName: "UserIntoRoom",
			Handler:    _Game_UserIntoRoom_Handler,
		},
		{
			MethodName: "SearchRoom",
			Handler:    _Game_SearchRoom_Handler,
		},
		{
			MethodName: "DeleteRoom",
			Handler:    _Game_DeleteRoom_Handler,
		},
		{
			MethodName: "QuitRoom",
			Handler:    _Game_QuitRoom_Handler,
		},
		{
			MethodName: "UpdateRoom",
			Handler:    _Game_UpdateRoom_Handler,
		},
		{
			MethodName: "UpdateUserReadyState",
			Handler:    _Game_UpdateUserReadyState_Handler,
		},
		{
			MethodName: "BeginGame",
			Handler:    _Game_BeginGame_Handler,
		},
		{
			MethodName: "BackRoom",
			Handler:    _Game_BackRoom_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "game.proto",
}
