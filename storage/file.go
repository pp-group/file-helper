package storage

import (
	"os"
	"path/filepath"
)

var _ IStorage = new(FileStorage)

type FileStorage struct {
	folder string
}

func NewFileStorage(folder string) (*FileStorage, error) {
	return &FileStorage{
		folder: folder,
	}, nil
}

func (storage *FileStorage) Writer(fileName string, help ParamsHelper) (IWriteBroker, error) {
	return NewLocalFileStorageBroker(storage.folder, fileName)
}
func (storage *FileStorage) Reader(fileName string, help ParamsHelper) (IReadBroker, error) {
	return NewLocalFileStorageBroker(storage.folder, fileName)
}
func (storage *FileStorage) Manager(fileName string, help ParamsHelper) (IManageBroker, error) {
	return NewLocalFileStorageBroker(storage.folder, fileName)
}

var _ IBroker = new(LocalFileStorageBroker)

type LocalFileStorageBroker struct {
	*os.File
	folder string
}

func NewLocalFileStorageBroker(folder string, fileName string) (*LocalFileStorageBroker, error) {

	var f *os.File

	{
		var err error
		if fileName != "" {
			f, err = os.OpenFile(filepath.Join(folder, fileName), os.O_CREATE|os.O_RDWR, os.ModePerm)
			if err != nil {
				return nil, err
			}

		}
	}

	return &LocalFileStorageBroker{
		File:   f,
		folder: folder,
	}, nil
}

func (broker *LocalFileStorageBroker) Close() error {
	return broker.File.Close()
}
func (broker *LocalFileStorageBroker) Exist(fullPath string) (bool, error) {
	file, err := os.Open(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	file.Close()
	return true, nil
}
func (broker *LocalFileStorageBroker) URL(filename string) (string, error) {
	fullPath := filepath.Join(broker.folder, filename)
	b, err := broker.Exist(fullPath)
	if err != nil {
		return "", nil
	}
	if !b {
		return "", os.ErrNotExist
	}
	return fullPath, nil

}
func (broker *LocalFileStorageBroker) Delete(filename string) error {
	return os.Remove(filepath.Join(broker.folder, filename))
}
func (broker *LocalFileStorageBroker) Read(p []byte) (n int, err error) {
	return broker.File.Read(p)
}
func (broker *LocalFileStorageBroker) Write(p []byte) (n int, err error) {
	return broker.File.Write(p)
}
