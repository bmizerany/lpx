package lpx_test

import (
	"bufio"
	"bytes"
	"github.com/bmizerany/lpx"
	"io"
	"log"
	"net/http"
)

func ExampleReader_Read() {
	// a simple HTTP server that echos router payloads
	h := func(w http.ResponseWriter, r *http.Request) {
		lr := lpx.NewReader(bufio.NewReader(r.Body))
		for lr.Next() {
			hdr := lr.Header()
			if bytes.Equal(hdr.Name, []byte("router")) {
				io.Copy(w, lr)
			}
		}
	}

	http.HandleFunc("/drain", h)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
