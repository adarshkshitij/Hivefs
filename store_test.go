package main

import (
	"bytes"
	"testing"
)
func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}
	s := NewStore(opts)
	data := bytes.NewReader([]byte("some jpg bytes"))
	err := s.writeStream("mypicture", data)
	if err != nil{
		t.Error(err)
	}
}