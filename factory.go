package factory

import "github.com/pp-group/file_helper/storage"

type StorageFactory func() storage.IStorage

func OssStorageFactory(endpoint, ak, sk string) func() (storage.IStorage, error) {
	return func() (storage.IStorage, error) {
		return storage.NewOssStorage(endpoint, ak, sk)
	}
}

func FileStorageFactory() func() (storage.IStorage, error) {
	return func() (storage.IStorage, error) {
		return storage.NewFileStorage()
	}
}
