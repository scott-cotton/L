package L

import (
	"io"
	"os"
	"runtime"
	"strings"
)

// Config is the configuration of a logger.
type Config struct {
	// Labels represents the set of labels associated with
	// a logger.
	Labels map[string]int

	// Pre is a sequence of Middlewares to pre-process
	// loggable objects.
	Pre []Middleware
	// Post is a sequence of Middlewares to post-process
	// loggable objects.  Post middlewares are called
	// after .Log() and before .F.Fmt().
	Post []Middleware
	// W is the writer for this logger.
	W io.Writer
	// F is a Fmter for the logger.
	F Fmter
	// E is a handler for any errors which occur during
	// formatting.
	E func(Logger, *Config, error)

	pkg string
}

// NewConfig returns a *Config with the associated labels.  The labels are
// subject to the following expansion rules
//
//  - a label starting with a colon ':' is prefixed with package name
func NewConfig(labels ...string) *Config {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc).Name()
	i := strings.LastIndexByte(fn, byte('.'))
	pkg := fn
	if i != -1 {
		pkg = fn[:i]
	}
	c := &Config{
		Labels: map[string]int{},
		W:      os.Stderr,
		E:      EPanic,
		pkg:    pkg,
	}
	for _, lbl := range labels {
		key := lbl
		if lbl != "" && lbl[0] == '.' {
			key = pkg + lbl
		}
		c.Labels[key] = 0
	}
	return c
}

// Clone clones the configuration c.
func (c *Config) Clone() *Config {
	res := &Config{}
	*res = *c
	res.Pre = append([]Middleware{}, c.Pre...)
	res.Post = append([]Middleware{}, c.Post...)
	res.Labels = make(map[string]int, len(c.Labels))
	res.pkg = c.pkg
	for k, v := range c.Labels {
		res.Labels[k] = v
	}
	return res
}

// Apply applies the configuration o to c.  Fields
// are copied over if they are not nil in o, otherwise
// left untouched.  Labels in o should not include
// the package name, but if they start with '.', they
// are expanded with the package name of 'c' when
// copied to c's Labels.
//
// if a label in 'c', with any package name stripped,
// is not in o, then it is removed from c.
func (c *Config) Apply(o *Config) {
	if o.E != nil {
		c.E = o.E
	}
	if o.W != nil {
		c.W = o.W
	}
	if o.F != nil {
		c.F = o.F
	}
	if o.Pre != nil {
		c.Pre = append([]Middleware{}, o.Pre...)
	}
	if o.Post != nil {
		c.Post = append([]Middleware{}, o.Post...)
	}
	if o.Labels == nil {
		return
	}
	if c.Labels == nil {
		c.Labels = make(map[string]int, len(o.Labels))
	}
	for k := range c.Labels {
		if _, ok := o.Labels[c.Unlocalize(k)]; !ok {
			delete(c.Labels, k)
		}
	}
	for k, v := range o.Labels {
		c.Labels[c.Localize(k)] = v
	}
}

func (c *Config) Unlocalize(label string) string {
	if label != "" && label[0] == '.' {
		return c.pkg + label
	}
	return label
}

func (c *Config) Localize(label string) string {
	if strings.HasPrefix(label, c.pkg+".") {
		return label[len(c.pkg)+1:]
	}
	return label
}

// CurrentRootConfig retrieves a clone of the configuration
// from the last call to Apply, if any.
func CurrentRootConfig() *Config {
	return root.ReadConfig()
}
