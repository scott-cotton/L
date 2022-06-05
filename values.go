package L

type Values interface {
	Dict() *Obj
	Array() *Obj
	Str(v string) *Obj
	Float(v float64) *Obj
	Bool(v bool) *Obj
	Int(v int) *Obj
	Bytes(d []byte) *Obj
	Null() *Obj
}
