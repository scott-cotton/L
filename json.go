package L

import (
	"io"
)

type jsonFmter struct {
}

func JSONFmter() Fmter {
	return jsonFmter{}
}

func (j jsonFmter) Fmt(w io.Writer, d []byte) error {
	_, e := w.Write(append(d, '\n'))
	return e
}
