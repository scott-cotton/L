package L

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type TableFmter struct {
	Fields    []string
	Sep       string
	Keys      bool
	FloatFmt  byte
	FloatPrec int
	buf       *bytes.Buffer
}

func (c *TableFmter) Fmt(w io.Writer, d []byte) error {
	if c.buf == nil {
		c.buf = bytes.NewBuffer(nil)
	}
	j := map[string]any{}
	json.Unmarshal(d, &j)
	enc := json.NewEncoder(c.buf)
	i := 0
	for _, f := range c.Fields {
		v := j[f]
		if v == nil {
			continue
		}
		if i != 0 {
			fmt.Fprintf(c.buf, c.Sep)
		}
		i++
		if c.Keys {
			fmt.Fprintf(c.buf, "%s=", f)
		}
		switch x := v.(type) {
		case string:
			fmt.Fprint(c.buf, x)
		case int:
			fmt.Fprint(c.buf, strconv.Itoa(x))
		case float64:
			fmt.Fprint(c.buf, strconv.FormatFloat(x, c.FloatFmt, c.FloatPrec, 64))
		case bool:
			fmt.Fprint(c.buf, strconv.FormatBool(x))
		default:
			if e := enc.Encode(v); e != nil {
				return e
			}
			// remove trailing newline
			c.buf.Truncate(c.buf.Len() - 1)
		}
	}
	fmt.Fprintln(c.buf)
	_, e := w.Write(c.buf.Bytes())
	return e
}
