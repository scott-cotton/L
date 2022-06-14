package L

import (
	"bytes"
	"context"
	"testing"
)

var ctxW = bytes.NewBuffer(nil)
var ctxCfg = &Config{
	W: ctxW,
	F: JSONFmter(),
}

func TestContext(t *testing.T) {
	ctx := context.TODO()
	l := New(ctxCfg)
	c := l.With(".warn", 1)
	ctx = Context(ctx, l.Dict().Field("a", 10))
	obj := FromContextWith(ctx, c)
	obj.Log()
	want := `{"a":10}` + "\n"
	if ctxW.String() != want {
		t.Errorf("got %q want %q", ctxW.String(), want)
	}
}
