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
	fileName string
	folder   string
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
		File:     f,
		fileName: filepath.Join(folder, fileName),
		folder:   folder,
	}, nil
}

func (broker *LocalFileStorageBroker) Close() error {
	return broker.File.Close()
}
func (broker *LocalFileStorageBroker) Exist() (bool, error) {
	file, err := os.Open(broker.fileName)
	if os.IsNotExist(err) {
		return false, nil
	}
	file.Close()
	return true, nil
}
func (broker *LocalFileStorageBroker) URL() (string, error) {
	b, err := broker.Exist()
	if err != nil {
		return "", nil
	}
	if !b {
		return "", os.ErrNotExist
	}
	return broker.fileName, nil

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
