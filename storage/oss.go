package storage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var _ IStorage = new(OssStorage)

type OssStorage struct {
	*oss.Client
	folder string
}

func NewOssStorage(endpoint, ak, sk, folder string) (*OssStorage, error) {
	client, err := oss.New(endpoint, ak, sk)
	if err != nil {
		return nil, err
	}
	return &OssStorage{
		Client: client,
		folder: folder,
	}, nil
}

func (storage *OssStorage) Writer(objFullPath string, helper ParamsHelper) (IWriteBroker, error) {
	broker, err := newOssStorageBroker(storage.Client, objFullPath, helper, storage.folder)

	if err != nil {
		return nil, err
	}

	return broker, nil
}

func (storage *OssStorage) Reader(objFullPath string, helper ParamsHelper) (IReadBroker, error) {

	broker, err := newOssStorageBroker(storage.Client, objFullPath, helper, storage.folder)

	if err != nil {
		return nil, err
	}
	return broker, nil
}

func (storage *OssStorage) Manager(objFullPath string, helper ParamsHelper) (IManageBroker, error) {
	broker, err := newOssStorageBroker(storage.Client, objFullPath, helper, storage.folder)

	if err != nil {
		return nil, err
	}
	return broker, nil
}

func newOssStorageBroker(client *oss.Client, objFullPath string, helper ParamsHelper, folder string) (*OssStorageBroker, error) {

	if helper == nil {
		return nil, errors.New("oss storage must need a paramsHelper to specify buckert name")
	}

	broker, err := NewOssStorageBroker(client, helper().(string), filepath.Join(folder, objFullPath))
	if err != nil {
		return nil, fmt.Errorf("initialize the broker err. %s", err.Error())
	}

	return broker, nil
}

var _ IBroker = new(OssStorageBroker)

type OssStorageBroker struct {
	IStorage
	bucket  *oss.Bucket
	nextPos int64
	objName string
	stream  io.ReadCloser
}

func NewOssStorageBroker(client *oss.Client, bucketName, objName string) (*OssStorageBroker, error) {
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	return &OssStorageBroker{
		nextPos: 0,
		bucket:  bucket,
		objName: objName,
	}, nil
}
func (broker *OssStorageBroker) download(objName string) (io.ReadCloser, error) {

	uploadStatus, err := broker.getObjMeta(broker.objName, "File-Upload-Status")

	if err != nil {
		return nil, err
	}

	if uploadStatus != "Finished" {
		return nil, errors.New("target file upload pending")
	}

	{
		var err error
		broker.stream, err = broker.bucket.GetObject(objName)
		if err != nil {
			return nil, err
		}
	}

	return broker.stream, nil

}

func (broker *OssStorageBroker) URL(objName string) (string, error) {

	uploadStatus, err := broker.getObjMeta(broker.objName, "File-Upload-Status")

	if err != nil {
		return "", err
	}

	if uploadStatus != "Finished" {
		return "", errors.New("target file upload pending")
	}

	signURL, err := broker.bucket.SignURL(objName, oss.HTTPGet, 60)

	if err != nil {
		return "", err
	}

	return signURL, nil

}

func (storage *OssStorageBroker) upload(objName string, objValue io.Reader) error {
	bucket := storage.bucket
	{
		var err error
		if storage.nextPos == 0 {
			storage.nextPos, err = bucket.AppendObject(objName, objValue, storage.nextPos, oss.Meta("File-Upload-Status", "Pending"))
		} else {
			storage.nextPos, err = bucket.AppendObject(objName, objValue, storage.nextPos)
		}
		if err != nil {
			return err
		}
	}
	storageType := oss.ObjectStorageClass(oss.StorageStandard)
	objectAcl := oss.ObjectACL(oss.ACLDefault)
	return bucket.PutObject(objName, objValue, storageType, objectAcl)
}

func (broker *OssStorageBroker) setObjMeta(objName, key, val string) error {
	return broker.bucket.SetObjectMeta(objName, oss.Meta(key, val))
}
func (broker *OssStorageBroker) getObjMeta(objName string, key string) (string, error) {

	h, err := broker.bucket.GetObjectMeta(objName)
	if err != nil {
		return "", err
	}
	return h.Get(oss.HTTPHeaderOssMetaPrefix + key), nil

}

func (broker *OssStorageBroker) Read(p []byte) (int, error) {

	var err error

	if broker.stream != nil {
		broker.stream, err = broker.download(broker.objName)
		if err != nil {
			return 0, err
		}
	}
	return broker.stream.Read(p)

}

func (broker *OssStorageBroker) Write(p []byte) (n int, err error) {
	if err = broker.upload(broker.objName, bytes.NewReader(p)); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (broker *OssStorageBroker) Delete(objKey string) error {
	return broker.bucket.DeleteObject(objKey)
}

func (broker *OssStorageBroker) Close() error {
	if broker.stream != nil {
		return broker.stream.Close()
	} else {
		return broker.setObjMeta(broker.objName, "File-Upload-Status", "Finished")
	}
}

func (broker *OssStorageBroker) Exist(objNmae string) (bool, error) {
	isExist, err := broker.bucket.IsObjectExist(objNmae)
	if err != nil {
		return false, err
	}
	return isExist, nil

}
