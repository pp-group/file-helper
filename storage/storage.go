package storage

import (
	"io"
)

type IStorage interface {
	Writer(string, ParamsHelper) (IWriteBroker, error)
	Reader(string, ParamsHelper) (IReadBroker, error)
	Manager(string, ParamsHelper) (IManageBroker, error)
}

type IReadBroker interface {
	io.Reader
	BrokerBundle
}

type IWriteBroker interface {
	io.Writer
	BrokerBundle
}

type IBroker interface {
	IReadBroker
	IWriteBroker
	IManageBroker
}

type BrokerBundle interface {
	Close() error
	Exist(string) (bool, error)
	URL(string) (string, error)
}

type IManageBroker interface {
	Delete(string) error
}

type ParamsHelper func() interface{}

var OssStorageBrokerParamHelperExample ParamsHelper = func() interface{} {
	return "your_bucket_naem"
}
