package store

import (
	"os"
	"testing"
)

func TestUpload(t *testing.T) {
	file, err := os.Open("file.go")
	if err != nil {
		t.Fatal(err)
	}
	path, err := Upload(file, "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
}
