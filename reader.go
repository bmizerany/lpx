package lpx

import (
	"bytes"
	"io"
	"strconv"
)

type BytesReader interface {
	io.Reader
	ReadBytes(delim byte) (line []byte, err error)
}

// A Header represents a single header in a logplex entry. All fields are
// popluated.
type Header struct {
	Time []byte
	Name []byte
}

// A Reader provides sequential access to logplex packages. The Next method
// advances to the next entry (including the first), and then can be treated
// as an io.Reader to access the packages payload.
type Reader struct {
	r   BytesReader
	b   io.Reader
	hdr *Header
	err error
	n   int64
}

// NewReader creates a new Reader reading from r.
func NewReader(r BytesReader) *Reader {
	return &Reader{r: r, hdr: new(Header)}
}

// Next advances to the next entry in the stream.
func (r *Reader) Next() bool {
	var l []byte
	r.field(&l) // message length
	r.n, r.err = strconv.ParseInt(string(l), 10, 64)
	if r.err != nil {
		return false
	}
	r.field(nil)         // PRI/VERSION
	r.field(&r.hdr.Time) // TIMESTAMP
	r.field(nil)         // HOSTNAME
	r.field(&r.hdr.Name) // APP-NAME
	r.field(nil)         // PROCID
	r.field(nil)         // MSGID
	r.b = io.LimitReader(r.r, r.n)
	return true
}

// Header returns the current entries decoded header.
func (r *Reader) Header() *Header {
	return r.hdr
}

// Read reads from the current entries payload. It returns 0, io.EOF when it
// reaches the end of the entries payload, until Next is called to advance the
// entry.
func (r *Reader) Read(b []byte) (n int, err error) {
	n, r.err = r.b.Read(b)
	return n, r.err
}

// Err returns the first non-EOF error that was encountered by the Reader.
func (r *Reader) Err() error {
	if r.err == io.EOF {
		return nil
	}
	return r.err
}

func (r *Reader) field(b *[]byte) {
	g, err := r.r.ReadBytes(' ')
	if err != nil {
		r.err = err
		return
	}
	r.n -= int64(len(g))
	if b == nil {
		return
	}
	*b = bytes.TrimRight(g, " ")
}
