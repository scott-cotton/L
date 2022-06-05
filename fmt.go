package L

import "io"

// a Fmter formats logs
type Fmter interface {
	Fmt(w io.Writer, d []byte) error
}
