package file_helper_test

import (
	"strings"
	"testing"

	"github.com/pp-group/file_helper"
)

func TestLocalBaseStorageFeature(t *testing.T) {
	factory := file_helper.FileStorageFactory("./temp")

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
