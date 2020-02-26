package storage

import (
	"github.com/thingio/edge-core/common/conf"
	"github.com/thingio/edge-core/common/proto/resource"
)

type ResourceStorage interface {
	Init() error
	Put(*resource.Resource) error
	Get(resource.Key) (*resource.Resource, error)
	List(resource.Key) ([]*resource.Resource, error)
}

func NewResourceStorage(config conf.DBConfig) ResourceStorage{
	return &BoltStorage{File:config.File, Timeout: config.Timeout}
}