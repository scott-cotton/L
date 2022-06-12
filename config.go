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
	Labels map[string]int `json:"labels,omitempty"`

	// Pre is a sequence of Middlewares to pre-process
	// loggable objects.
	Pre []Middleware `json:"-"`

	// Post is a sequence of Middlewares to post-process
	// loggable objects.  Post middlewares are called
	// after .Log() and before .F.Fmt().
	Post []Middleware `json:"-"`

	// W is the writer for this logger.
	W io.Writer `json:"-"`

	// F is a Fmter for the logger.
	F Fmter `json:"-"`

	// E is a handler for any errors which occur during
	// formatting.
	E func(*Config, error) `json:"-"`

	pkg string
}

// NewConfig returns a *Config with the associated labels.  The labels are
// subject to the following expansion rules
//
//  - a label starting with '.' is prefixed with package name.
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
		c.Labels[c.Unlocalize(lbl)] = 0
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

// Unlocalize prefixes label with the c's package if
// label starts with '.'.
func (c *Config) Unlocalize(label string) string {
	if label != "" && label[0] == '.' {
		return c.pkg + label
	}
	return label
}

// localize strips the prefix '<pkg>.' if label
// starts with c's package.
func (c *Config) Localize(label string) string {
	if strings.HasPrefix(label, c.pkg+".") {
		return label[len(c.pkg):]
	}
	return label
}

// CurrentRootConfig retrieves a clone of the configuration
// from the last call to Apply, if any.
func CurrentRootConfig() *Config {
	return root.ReadConfig()
}

func (c *Config) Package() string {
	return c.pkg
}

type PackageConfig struct {
	Config
	Package string `json:"package"`
}

type ConfigNode struct {
	PackageConfig
	Parent int `json:"parent"`
}
