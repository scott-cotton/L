package L

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type TableFmter struct {
	Fields []string
	Sep    string
	buf    *bytes.Buffer
}

func (c *TableFmter) Fmt(w io.Writer, d []byte) error {
	if c.buf == nil {
		c.buf = bytes.NewBuffer(nil)
	}
	j := map[string]any{}
	json.Unmarshal(d, &j)
	enc := json.NewEncoder(c.buf)
	for _, f := range c.Fields {
		v := j[f]
		if v == nil {
			continue
		}
		fmt.Fprintf(c.buf, "%s=", f)
		if e := enc.Encode(v); e != nil {
			return e
		}
		// remove trailing newline
		c.buf.Truncate(c.buf.Len() - 1)
		fmt.Fprint(c.buf, c.Sep)
	}
	fmt.Fprintln(c.buf)
	_, e := w.Write(c.buf.Bytes())
	return e
}
