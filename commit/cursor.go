package commit

// A Cursor can iterate over a sequence of commits.
type Cursor struct {
	seq   Sequence
	start SeqNum
	end   SeqNum
	pos   SeqNum
}

// NewCursor returns a commit cursor for r. The cursor will iterate
// over a sequence of commits, up to the most recent commit in r
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

// First moves the cursor to the first commit in the sequence.
func (c *Cursor) First() {
	c.pos = c.start
}

// Last moves the cursor to the last commit in the sequence.
func (c *Cursor) Last() {
	c.pos = c.end - 1
}

// Next moves the cursor to the next commit in the sequence.
func (c *Cursor) Next() {
	c.pos++
}

// Previous moves the cursor to the previous commit in the sequence.
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
