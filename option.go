package drivestream

import "io"

// Option is a configuration option for a stream.
type Option func(*Stream)

// WithLogger causes the stream to write log output to w.
func WithLogger(w io.Writer) Option {
	return func(s *Stream) {
		s.stdout = w
	}
}
