package utils

import (
	"fmt"
	"store_srv/global"

	"github.com/go-redsync/redsync/v4"
)

func GetRedSync(roomID uint32) (*redsync.Mutex, error) {
	mutex := global.RedSync.NewMutex(fmt.Sprintf("room%d_lock", roomID))
	if err := mutex.Lock(); err != nil {
		return nil, err
	}
	return mutex, nil
}

func ReleaseRedSync(mutex *redsync.Mutex) error {
	if ok, err := mutex.Unlock(); !ok || err != nil {
		return err
	}
	return nil
}
