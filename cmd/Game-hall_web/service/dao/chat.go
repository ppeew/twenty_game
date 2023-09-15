package dao

import (
	"context"
	"hall_web/service/domains"
)

type ChatDao struct {
}

func (cd *ChatDao) GetChatList(ctx context.Context) ([]*domains.Message, error) {
	return nil, nil
}

func (cd *ChatDao) PushChat(ctx context.Context, message *domains.Message) error {
	return nil
}

func (cd *ChatDao) PopChat(ctx context.Context) error {
	return nil
}

func (cd *ChatDao) ChangeChat(ctx context.Context, i int) error {
	return nil
}
