package L

import "context"

var objKey = struct{}{}

// Context creates a new child context of 'p'
func Context(p context.Context, obj *Obj) context.Context {
	return context.WithValue(p, objKey, obj)
}

// FromContext retrieves a clone of the object from ctx if any.
func FromContext(ctx context.Context) *Obj {
	v := ctx.Value(objKey)
	if v == nil {
		return nil
	}
	o := v.(*Obj)
	return o.Clone()
}

// FromContextWith retrieves a clone of the object from
// ctx, if any.  The result is guaranteed to be associated
// with 'l' if not nil.
func FromContextWith(ctx context.Context, l Logger) *Obj {
	o := FromContext(ctx)
	if o == nil {
		return nil
	}
	r := o.getRoot()
	r.logger = l
	return r
}
