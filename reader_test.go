package lpx

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestReader(t *testing.T) {
	const data = `66 <174>1 2012-07-22T00:06:26-00:00 somehost Go console 2 Hi from Go
67 <174>1 2013-07-22T00:06:26-00:00 somehost Go console 10 Hi from Py
`
	r := NewReader(bufio.NewReader(bytes.NewBufferString(data)))

	if !r.Next() {
		t.Error("want next")
	}

	w := &Header{[]byte("2012-07-22T00:06:26-00:00"), []byte("Go")}
	if !reflect.DeepEqual(w, r.Header()) {
		t.Errorf("want %q, got %q", w, r.Header())
	}

	if g := string(r.Bytes()); g != "Hi from Go\n" {
		t.Errorf("want %q, got %q", g, "Hi from Go\n")
	}
	if r.Err() != nil {
		t.Errorf("want %v, got %v", nil, r.Err())
	}

	if !r.Next() {
		t.Error("want next")
	}

	w = &Header{[]byte("2013-07-22T00:06:26-00:00"), []byte("Go")}
	if !reflect.DeepEqual(w, r.Header()) {
		t.Errorf("want %q, got %q", w, r.Header())
	}

	if g := string(r.Bytes()); g != "Hi from Py\n" {
		t.Errorf("want %q, got %q", g, "Hi from Py\n")
	}
	if r.Err() != nil {
		t.Errorf("want %v, got %v", nil, r.Err())
	}

	if r.Next() {
		t.Error("did not want next")
	}
}
