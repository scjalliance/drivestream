package binpath

// Text is a slice of strings that can be marshaled as a binary path.
type Text []string

// String returns p marshaled as a string.
func (p Text) String() string {
	buf := make([]byte, p.EncodedLen())
	p.MarshalTextBuffered(buf)
	return string(buf)
}

// Bytes returns p marshaled as a slice of bytes.
func (p Text) Bytes() (text []byte) {
	buf := make([]byte, p.EncodedLen())
	return p.MarshalTextBuffered(buf)
}

// MarshalText returns p marshaled as text. It never returns an error.
func (p Text) MarshalText() (text []byte, err error) {
	return p.Bytes(), nil
}

// MarshalTextBuffered marshals p into buf and returns the slice of bytes
// that were marshaled. If buf is nil or of insufficient length, a buffer of
// sufficient size will be allocated and returned.
func (p Text) MarshalTextBuffered(buf []byte) []byte {
	needed := p.EncodedLen()
	if needed == 0 {
		return nil
	}
	if len(buf) < needed {
		buf = make([]byte, needed)
	}
	i := copy(buf, []byte(p[0]))
	for c := 1; c < len(p); c++ {
		buf[i] = '/'
		i++
		i += copy(buf[i:], []byte(p[c]))
	}
	return buf[:needed]
}

// EncodedLen returns the encoded length of p.
func (p Text) EncodedLen() int {
	if len(p) == 0 {
		return 0
	}
	length := len(p) - 1 // Slashes
	for _, component := range p {
		length += len(component)
	}
	return length
}
