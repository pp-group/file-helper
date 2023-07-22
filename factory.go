package file_helper

import "github.com/pp-group/file_helper/storage"

type StorageFactory func() storage.IStorage

func OssStorageFactory(endpoint, ak, sk, folder string) func() (storage.IStorage, error) {
	return func() (storage.IStorage, error) {
		return storage.NewOssStorage(endpoint, ak, sk, folder)
	}
}

func FileStorageFactory(folder string) func() (storage.IStorage, error) {
	return func() (storage.IStorage, error) {
		return storage.NewFileStorage(folder)
	}
}
