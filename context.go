package L

import "context"

var logKey = struct{}{}
var objKey = struct{}{}

func Context(ctx context.Context, obj *Obj) context.Context {
	return context.WithValue(ctx, objKey, obj)
}

func FromContext(ctx context.Context, key string) *Obj {
	v := ctx.Value(logKey)
	o := v.(*Obj)
	return o.Clone()
}
