package file_helper_test

import (
	"os"
	"strings"
	"testing"

	file_helper "github.com/pp-group/file-helper"
)

//todo add the testcase

const (
	OssEndpoint = "OSS_ENDPOINT"
	OssAK       = "OSS_AK"
	OssSK       = "OSS_SK"
)

func OssConfig() (endpoint string, ak string, sk string) {
	return os.Getenv(OssEndpoint), os.Getenv(OssAK), os.Getenv(OssSK)
}

func TestOssBaseStorageFeature(t *testing.T) {

	endpoint, ak, sk := OssConfig()

	t.Logf("endpoint: %s , ak: %s, sk: %s", endpoint, ak, sk)

	factory := file_helper.OssStorageFactory(endpoint, ak, sk, "default")

	ossStorage, err := factory()
	if err != nil {
		t.Fatal(err)
	}

	file_name := "example.txt"

	writerBroker, err := ossStorage.Writer(file_name, func() interface{} {
		return "aigc-tts"
	})

	if err != nil {
		t.Fatal(err)
	}

	var text []string = []string{"{`csh0101`}"}

	n, err := writerBroker.Write([]byte(strings.Join(text, "\t")))

	if err != nil {
		t.Fatal(err)
	}

	t.Log("write bytes len", n)

	if err := writerBroker.Close(); err != nil {
		t.Fatal(err)
	}

	readerBroker, err := ossStorage.Reader(file_name, func() interface{} {
		return "aigc-tts"
	})

	if err != nil {
		t.Fatal(err)
	}

	fileURL, err := readerBroker.URL()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(fileURL)

	// manangeBroker, err := ossStorage.Manager(file_name, func() interface{} {
	// 	return "aigc-tts"
	// })

	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if err := manangeBroker.Delete(file_name); err != nil {
	// 	t.Fatal(err)
	// }
}
