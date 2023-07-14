package handler

import (
	"context"
	"store_srv/proto/store"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StoreServer struct {
	store.UnimplementedStoreServer
}

func (StoreServer) SelectGoods(context.Context, *store.SelectGoodsReq) (*store.SelectGoodsRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SelectGoods not implemented")
}
func (StoreServer) BuyGoods(context.Context, *store.BuyGoodReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuyGoods not implemented")
}
