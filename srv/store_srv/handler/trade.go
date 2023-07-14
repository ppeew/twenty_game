package handler

import (
	"context"
	"store_srv/proto/store"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (StoreServer) SelectTradeItems(context.Context, *store.SelectTradeItemsReq) (*store.SelectTradeItemsRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SelectTradeItems not implemented")
}
func (StoreServer) PushTradeItem(context.Context, *store.PushTradeItemReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushTradeItem not implemented")
}
func (StoreServer) DownTradeItem(context.Context, *store.DownTradeItemReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownTradeItem not implemented")
}
func (StoreServer) BuyTradeItem(context.Context, *store.BuyTradeItemReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuyTradeItem not implemented")
}
