package state

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/ngyewch/expframework/codec"
	"google.golang.org/protobuf/proto"
)

type Manager[T any, PT interface {
	*T
	proto.Message
}] struct {
	appId   string
	groupId string
	codec   codec.Codec[T]
}

func NewManager[T any, PT interface {
	*T
	proto.Message
}](appId string, groupId string, codec codec.Codec[T]) *Manager[T, PT] {
	return &Manager[T, PT]{
		appId:   appId,
		groupId: groupId,
		codec:   codec,
	}
}

func (manager *Manager[T, PT]) LoadState(id string) (*T, error) {
	stateFile, err := xdg.StateFile(filepath.Join(manager.appId, manager.groupId, id+".json"))
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	state, err := manager.codec.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (manager *Manager[T, PT]) SaveState(id string, state *T) error {
	stateFile, err := xdg.StateFile(filepath.Join(manager.appId, manager.groupId, id+".json"))
	if err != nil {
		return err
	}
	b, err := manager.codec.MarshalBinary(state)
	if err != nil {
		return err
	}
	err = os.WriteFile(stateFile, b, 0755)
	if err != nil {
		return err
	}
	return nil
}
