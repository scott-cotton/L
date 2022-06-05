package L

import (
	"runtime"
	"strings"
)

// Middleware is a type for hooks into Loggers'
// configuration.
type Middleware func(Logger, *Config, *Obj) *Obj

// Pkg() is a Middleware which injects a field
// with key "Lpkg" and value of the package path
// of the config.
func Pkg() Middleware {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc).Name()
	i := strings.LastIndexByte(fn, byte('.'))
	pkg := fn
	if i != -1 {
		pkg = fn[:i]
	}
	return func(_ Logger, _ *Config, o *Obj) *Obj {
		return o.Field("Lpkg", pkg)
	}
}

// If is a middleware that filters objects when
// the config does not contain the label 'label'.
//
// If can be used as a Pre middleware, in which
// case the overhead of message construction is
// eliminated in addition to the overhead of
// message writing.
func If(label string) Middleware {
	return func(ell Logger, cfg *Config, o *Obj) *Obj {
		if _, present := cfg.Labels[label]; present {
			return o
		}
		return nil
	}
}

func IfNot(label string) Middleware {
	return func(ell Logger, cfg *Config, o *Obj) *Obj {
		if _, present := cfg.Labels[label]; present {
			return nil
		}
		return o
	}
}

func Leq(label string, value int) Middleware {
	return func(ell Logger, cfg *Config, o *Obj) *Obj {
		if cfg.Labels[label] <= value {
			return o
		}
		return nil
	}
}

func Geq(label string, value int) Middleware {
	return func(ell Logger, cfg *Config, o *Obj) *Obj {
		if cfg.Labels[label] >= value {
			return o
		}
		return nil
	}
}

func Label(label string) Middleware {
	return func(ell Logger, cfg *Config, o *Obj) *Obj {
		return o.Field(label, cfg.Labels[label])
	}
}
