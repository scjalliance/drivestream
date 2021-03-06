package collection

// A Cursor can iterate over a sequence of collections.
type Cursor struct {
	seq   Sequence
	start SeqNum
	end   SeqNum
	pos   SeqNum
}

// NewCursor returns a collection cursor for seq. The cursor will iterate
// over the sequence of collections, up to the most recent collection in seq
// at the time the cursor was created.
func NewCursor(seq Sequence) (*Cursor, error) {
	end, err := seq.Next()
	if err != nil {
		return nil, err
	}
	return &Cursor{seq: seq, end: end}, nil
}

// Valid return true if the current sequence number is valid.
func (c *Cursor) Valid() bool {
	return c.pos >= c.start && c.pos < c.end
}

// SeqNum returns the current sequence number.
func (c *Cursor) SeqNum() SeqNum {
	return c.pos
}

// First moves the cursor to the first collection in the sequence.
func (c *Cursor) First() {
	c.pos = c.start
}

// Last moves the cursor to the last collection in the sequence.
func (c *Cursor) Last() {
	c.pos = c.end - 1
}

// Next moves the cursor to the next collection in the sequence.
func (c *Cursor) Next() {
	c.pos++
}

// Previous moves the cursor to the previous collection in the sequence.
func (c *Cursor) Previous() {
	c.pos--
}

// Seek moves the cursor to the given sequence number.
func (c *Cursor) Seek(pos SeqNum) {
	c.pos = pos
}

// Reader returns a reader for the current sequence number.
func (c *Cursor) Reader() (*Reader, error) {
	return NewReader(c.seq.Ref(c.pos))
}

// Writer returns a writer for the current sequence number.
func (c *Cursor) Writer(instance string) (*Writer, error) {
	return NewWriter(c.seq.Ref(c.pos), instance)
}
