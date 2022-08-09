package test

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cast"
)

const chunkSize = 1024 * 1024 * 10

func TestGenerateChunkFile(t *testing.T) {
	fileInfo, err := os.Stat("test.mp4")
	if err != nil {
		t.Fatal(err)
	}
	chunkNum := fileInfo.Size() / chunkSize
	if fileInfo.Size()%chunkSize != 0 {
		chunkNum++
	}
	file, err := os.OpenFile("test.mp4", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	var b []byte
	for i := 0; i < int(chunkNum); i++ {
		b = make([]byte, min(chunkSize, int(fileInfo.Size())-i*chunkSize))
		file.Seek(int64(i*chunkSize), 0)
		file.Read(b)
		f, err := os.OpenFile("./"+cast.ToString(i)+".chunk", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		f.Write(b)
	}
}

func TestMergeChunkFile(t *testing.T) {
	file, err := os.OpenFile("test2.mp4", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	fileInfo, err := os.Stat("test.mp4")
	if err != nil {
		t.Fatal(err)
	}
	chunkNum := fileInfo.Size() / chunkSize
	if fileInfo.Size()%chunkSize != 0 {
		chunkNum++
	}
	for i := 0; i < int(chunkNum); i++ {
		b, err := ioutil.ReadFile("./" + cast.ToString(i) + ".chunk")
		if err != nil {
			t.Fatal(err)
		}
		file.Write(b)
	}
}

func TestCheck(t *testing.T) {
	f1, err := ioutil.ReadFile("test.mp4")
	if err != nil {
		t.Fatal(err)
	}
	f2, err := ioutil.ReadFile("test2.mp4")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%x", md5.Sum(f1)) == fmt.Sprintf("%x", md5.Sum(f2)) {
		fmt.Println("校验通过")
	} else {
		t.Fatal("校验失败")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
