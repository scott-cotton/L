package L

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf8"
)

// Obj represents a thing to be logged, with structured logging in mind.
// 'Obj' provides a set of chainable methods for constructing structured data
// and a small set of 'bytes.Buffer'-like methods.
//
// 'Objs' are not safe to use in multiple goroutines.
//
// When Objs are chained, all method calls occur by side-effect; the object
// is updated destructively.  However, the methods '.Dict' and '.Array' create
// and return new child 'Obj's while updating the parent.
//
// When an Obj is chained in this way, each parent should have only 1 child
// at a given time.  If one wishes to create a structure with multiple children,
//
//    c_1 := obj.Dict()
//    c_2 := obj.Dict()
//
// 'c_i' should be '.Close'd before 'c_{i+1}' is created.  If this is not
// the case, then the call will panic. One can use obj.Set("key", c2)
// with a fresh object if desired.
//
// This limitation is deliberate and aligned with logging, where
// there is usually a linear order to what is being written to a log.
//
// The nil '* Obj' will never panic.
//
// '* Obj' is an 'io.WriteCloser' to support direct encoding in various formats.
type Obj struct {
	parent   *Obj
	root     *Obj
	hasChild bool
	hadChild bool
	d        []byte
	i        int
	logger   *logger
}

func (t *Obj) mkChild() *Obj {
	if t.hasChild == true {
		return nil
	}
	res := &Obj{parent: t, i: len(*t.buf())}
	t.hasChild = true
	return res
}

func (t *Obj) getRoot() *Obj {
	if t.root != nil {
		return t.root
	}
	if t.parent == nil {
		t.root = t
		return t
	}
	t.root = t.parent.getRoot()
	return t.root
}

func (t *Obj) buf() *[]byte {
	return &t.getRoot().d
}

// Clone creates a clone of 't' which can be manipulated
// independently of 't'.
func (t *Obj) Clone() *Obj {
	if t == nil {
		return nil
	}
	var p *Obj
	if t.parent != nil {
		p = t.parent.Clone()
	}
	res := &Obj{parent: p, i: t.i, hasChild: t.hasChild}
	res.root = res.getRoot()
	d := *t.buf()

	if res.IsRoot() {
		dd := make([]byte, len(d))
		copy(dd, d)
		res.d = dd
	}
	return res
}

func (t *Obj) Dict() *Obj {
	if t == nil {
		return nil
	}
	if t.hasChild {
		panic("nonlinear")
	}
	if t.hadChild {
		t.WriteByte(',')
	}
	c := t.mkChild()
	c.WriteByte('{')
	t.hadChild = true
	return c
}

func (t *Obj) Array() *Obj {
	if t == nil {
		return nil
	}
	if t.hasChild {
		panic("nonlinear")
	}
	if t.hadChild {
		t.WriteByte(',')
	}
	c := t.mkChild()
	c.WriteByte('[')
	t.hadChild = true
	return c
}

// Err calls t.Field("Lerr", e.Error())
func (t *Obj) Err(e error) *Obj {
	return t.Field("Lerr", e.Error())
}

// Field sets a field with key 's' to value 'v'.
// 'v' must be a bool, int, string, float64, or
// *Obj or implement json.Marshaler.  If
// 'v' implements json.Marshaler, and json.Marshal(v)
// returns an error, Field panics.
func (t *Obj) Field(s string, v any) *Obj {
	if t == nil {
		return nil
	}
	if t.hadChild {
		t.WriteByte(',')
	}
	t.hadChild = false
	t = t.Str(s).WriteByte(':')
	t.hadChild = false
	switch x := v.(type) {
	case bool:
		t.Bool(x)
	case int:
		t.Int(x)
	case string:
		t.Str(x)
	case float64:
		t.Float(x)
	case []byte:
		t.Bytes(x)
	case *Obj:
		d := t.buf()
		*d = append(*d, *x.buf()...)
	default:
		if jm, ok := v.(json.Marshaler); ok {
			m, err := jm.MarshalJSON()
			if err != nil {
				panic(err)
			}
			d := t.buf()
			*d = append(*d, m...)
		}
	}
	t.hadChild = true
	return t
}

func (t *Obj) Str(s string) *Obj {
	if t == nil {
		return nil
	}
	if t.hadChild {
		t.WriteByte(',')
	}
	d := *t.buf()
	d = append(d, '"')
	for _, r := range s {
		if r == utf8.RuneError {
			d = append(d, []byte(`\ufffd`)...)
			continue
		}
		switch r {
		case '"', '\\':
			d = append(d, '\\', byte(r))
		case '\b':
			d = append(d, '\\', 'b')
		case '\f':
			d = append(d, '\\', 'f')
		case '\n':
			d = append(d, '\\', 'n')
		case '\r':
			d = append(d, '\\', 'r')
		default:
			d = utf8.AppendRune(d, r)
		}
	}
	d = append(d, '"')
	*t.buf() = d
	t.hadChild = true
	return t
}

func (t *Obj) Bool(v bool) *Obj {
	if t == nil {
		return nil
	}
	if t.hadChild {
		t.WriteByte(',')
	}
	if v {
		t.Write([]byte("true"))
		return t
	}
	t.Write([]byte("false"))
	t.hadChild = true
	return t
}

func (t *Obj) Int(i int) *Obj {
	if t == nil {
		return nil
	}
	if t.hadChild {
		t.WriteByte(',')
	}
	d := t.buf()
	*d = strconv.AppendInt(*d, int64(i), 10)
	t.hadChild = true
	return t
}

func (t *Obj) Float(v float64) *Obj {
	if t == nil {
		return nil
	}
	if t.hadChild {
		t.WriteByte(',')
	}
	d := t.buf()
	*d = strconv.AppendFloat(*d, v, 'e', -1, 64)
	t.hadChild = true
	return t
}

func Float(v float64) *Obj {
	res := &Obj{}
	return res.Float(v)
}

// Bytes encodes 'd' as a base64 encoded string.
func (t *Obj) Bytes(d []byte) *Obj {
	if t == nil {
		return nil
	}
	if t.hadChild {
		t.WriteByte(',')
	}
	t.WriteByte('"')
	r := t.buf()
	encLen := base64.StdEncoding.EncodedLen(len(d))
	n := len(*r)
	m := n + encLen
	if cap(*r) < m {
		tmp := make([]byte, m, m+n*2)
		*r = tmp
	}
	base64.StdEncoding.Encode(*r, d)
	t.WriteByte('"')
	t.hadChild = true
	return t
}

func (t *Obj) IsDict() bool {
	if t == nil {
		return false
	}
	return t.D()[t.i] == '{'
}

// Log logs the object with the logger that created it.
func (t *Obj) Log() {
	if t == nil {
		return
	}
	r := t.getRoot()
	r.logger.Log(t)
}

// D returns the underlying []byte.
func (t *Obj) D() []byte {
	if t == nil {
		return nil
	}
	return *t.buf()
}

// WriteByte writes b to the underlying []byte.
func (t *Obj) WriteByte(b byte) *Obj {
	if t == nil {
		return nil
	}
	r := t.buf()
	*r = append(*r, b)
	return t
}

// Null writes 'null' to the underlying []byte.
func (t *Obj) Null() *Obj {
	if t == nil {
		return nil
	}
	r := t.buf()
	*r = append(*r, []byte("null")...)
	t.hadChild = true
	return t
}

func Null() *Obj {
	res := &Obj{}
	return res.Null()
}

// Write implements io.Writer
func (t *Obj) Write(d []byte) (int, error) {
	if t == nil {
		return 0, nil
	}
	r := *t.buf()
	r = append(r, d...)
	*t.buf() = r
	return len(d), nil
}

// json encodes jm into t.
func (t *Obj) JSON(jm json.Marshaler) error {
	if t == nil {
		return nil
	}
	r := t.buf()
	buf := bytes.NewBuffer(*r)
	err := json.NewEncoder(buf).Encode(jm)
	if err != nil {
		return err
	}
	*r = buf.Bytes()
	t.hadChild = true
	return nil
}

func (t *Obj) IsRoot() bool {
	if t == nil {
		return false
	}
	return t.getRoot() == t
}

func (t *Obj) Parent() *Obj {
	if t == nil {
		return nil
	}
	return t.parent
}

// Close is a noop unless 't' was created with 'Dict' or 'Array' or
// 'Logger.Dict'.  In these cases, close validates the json output,
// returning any error.
//
// Note that Close is called by 'Logger.Log', so it is only necessary
// on roots created by calls to 'Dict' or 'Array'.
func (t *Obj) Close() error {
	if t == nil {
		return nil
	}
	d := *t.buf()
	switch d[t.i] {
	case '{':
		d = append(d, '}')
	case '[':
		d = append(d, ']')
	default:
		if t.parent != nil {
			return fmt.Errorf("invalid buffer")
		}
	}
	*t.buf() = d
	if !json.Valid(d) {
		var v any
		return json.Unmarshal(*t.buf(), &v)
	}
	return nil
}

func (t *Obj) _() error {
	d := *t.buf()
	switch d[t.i] {
	case '{':
		d = append(d, '}')
	case '[':
		d = append(d, ']')
	default:
		return fmt.Errorf("not a parent")
	}
	*t.buf() = d
	n, e := ck(d, t.i)
	if e != nil {
		return e
	}
	if n != len(d) {
		return fmt.Errorf("trailing characters at offset %d", n)
	}
	return nil
}

func ck(d []byte, i int) (n int, err error) {
	if len(d) == i {
		return i, nil
	}
	fmt.Printf("check %s %d\n", string(d), i)
	c := d[i]
	switch d[i] {
	case '{':
		n, err = ckObj(d, i+1)
	case '[':
		n, err = ckArray(d, i+1)
	case '"':
		n, err = ckStr(d, i+1)
	case 'n':
		n, err = ckNull(d, i+1)
	case 't':
		n, err = ckTrue(d, i+1)
	case 'f':
		n, err = ckFalse(d, i+1)
	default:
		if c >= '0' && c <= '9' || c == '-' {
			n, err = ckNumber(d, i)
			return
		}
		return i, fmt.Errorf("unexpected %c at offset %d", d[i], i)
	}
	return
}

func ckObj(d []byte, i int) (j int, err error) {
	switch d[i] {
	case '"':
	key:
		j, err = ckStr(d, i+1)
		if err != nil {
			return 0, err
		}
		if j == len(d) {
			return 0, fmt.Errorf("expected key, got end at %d", j)
		}
		if d[j] != ':' {
			return 0, fmt.Errorf("expected ':' at %d, got %c", j, d[j])
		}
		j, err = ck(d, j)
		if err != nil {
			return 0, err
		}
		if j == len(d) {
			return 0, fmt.Errorf("expected '}' or ',' at %d", j)
		}
		switch d[j] {
		case ',':
			i = j + 1
			if i == len(d) {
				return 0, fmt.Errorf("unexpected end")
			}
			if d[i] == '"' {
				goto key
			}
		case '}':
			return j + 1, nil
		}
	case '}':
		return i + 1, nil
	default:
		return 0, fmt.Errorf("unexpected character at %d", i)
	}
	return
}
func ckArray(d []byte, i int) (j int, err error) {
	if d[i] == ']' {
		j = i
		return
	}
value:
	j, err = ck(d, i)
	if err != nil {
		return 0, err
	}
	if j == len(d) {
		return 0, fmt.Errorf("unexpected end")
	}
	switch d[j] {
	case ',':
		i = j + 1
		goto value
	case ']':
		j++
		return
	default:
		return 0, fmt.Errorf("unexpected %c at %d", d[j], j)
	}
	return
}
func ckStr(d []byte, i int) (j int, err error) {
	fmt.Printf("check str at %s", d[i:])
	for k, b := range d[i:] {
		if b == '"' {
			return i + k + 1, nil
		}
	}
	return i, fmt.Errorf("unclosed '\"' at %d", i)
}
func ckNumber(d []byte, i int) (j int, err error) {
	if d[i] == '-' {
		i++
		if i == len(d) {
			return 0, fmt.Errorf("unexpected end")
		}
	}
	j = i
	for j < len(d) {
		switch {
		case d[j] >= '0' && d[j] <= '9':
		case d[j] == '.' || d[j] == 'e':
		default:
			return 0, fmt.Errorf("unexpected %c at %d", d[j], j)
		}
		j++
	}

	return
}
func ckTrue(d []byte, i int) (j int, err error) {
	if bytes.HasPrefix(d[i:], []byte("rue")) {
		return i + 3, nil
	}
	return i, fmt.Errorf("expected true at %d", i)
}
func ckFalse(d []byte, i int) (j int, err error) {
	if bytes.HasPrefix(d[i:], []byte("alse")) {
		return i + 4, nil
	}
	return i, fmt.Errorf("expected false at %d", i)
}
func ckNull(d []byte, i int) (j int, err error) {
	if bytes.HasPrefix(d[i:], []byte("ull")) {
		return i + 3, nil
	}
	return i, fmt.Errorf("expected null at %d", i)
}
