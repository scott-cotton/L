package L

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// Logger is the interface to a structured logger.
type Logger interface {
	// logs a structured object as JSON
	Log(o *Obj)

	Values

	// ReadConfig returns a clone of the Config associated with this
	// logger.
	ReadConfig() *Config

	// ApplyConfig calls 'lc.Apply(cfg)' on this logger and all of its
	// children where 'lc' is the config associated with this logger, or in
	// recursive application, the child logger.
	ApplyConfig(cfg *Config)

	// With returns a child logger with the same configuration
	// as this logger but which additionally has the specified
	// labels.  As in 'NewConfig', labels starting with a '.'
	// are implicitly package scoped to the package specified
	// in this loggers configuration.
	WithMap(map[string]int) Logger

	// With is convenience for WithMap(map[string]int{lbl: v}).
	With(lbl string, v int) Logger

	// Close closes this logger.  A global logger in an application need
	// not be closed.  However, any logger which is not global should be
	// closed or risk leaking underlying resources.
	Close() error
}

type logger struct {
	mu       sync.Mutex
	parent   *logger
	config   *Config
	children map[*logger]struct{}
	i        int
}

func (l *logger) With(key string, v int) Logger {
	return l.WithMap(map[string]int{key: v})
}

func (l *logger) WithMap(labels map[string]int) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	cfg := l.config.Clone()
	for lbl, v := range labels {
		key := lbl
		if lbl != "" && lbl[0] == '.' {
			key = cfg.pkg + lbl
		}
		cfg.Labels[key] = v
	}
	return New(cfg)
}

func (l *logger) ReadConfig() *Config {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.config.Clone()
}

// Log closes 'o' and if that results in an error 'e', it calls
// 'config.E(l, e)' where 'config' is the current configuration
// of 'l'.
func (l *logger) Log(o *Obj) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, mw := range l.config.Post {
		o = mw(l.config, o)
	}
	if err := o.Close(); err != nil {
		if l.config.E != nil {
			l.config.E(l.config, err)
		}
		return
	}
	if l.config.F != nil {
		l.config.F.Fmt(l.config.W, o.D())
	}
}

func (l *logger) ApplyConfig(cfg *Config) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.Apply(cfg)
	for c := range l.children {
		c.ApplyConfig(cfg)
	}
}

// New creates a new logger with a clone of 'cfg'.
func New(cfg *Config) Logger {
	cc := cfg.Clone()
	res := root.mkChild()
	res.config = cc
	if cfg.pkg == "" {
		pc, _, _, _ := runtime.Caller(1)
		fn := runtime.FuncForPC(pc).Name()
		i := strings.LastIndexByte(fn, byte('.'))
		pkg := fn
		if i != -1 {
			pkg = fn[:i]
		}
		cfg.pkg = pkg
	}
	return res
}

var root = &logger{
	config: &Config{
		W: os.Stderr,
		F: nil,
		E: EPanic,
	},
}

// Apply applies the configuration `c` to all Loggers created by this package.
func ApplyConfig(c *Config) {
	root.ApplyConfig(c)
}

func (l *logger) mkChild() *logger {
	res := &logger{
		parent: l,
		config: l.config,
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.children == nil {
		l.children = map[*logger]struct{}{}
	}
	l.children[res] = struct{}{}
	return res
}

func (l *logger) unlink(c *logger) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.children, l)
}

// Close closes the logger.  close does not close the associated
// Writer in the config.
func (l *logger) Close() error {
	if l == nil {
		return nil
	}
	if l.parent == nil {
		return nil
	}
	l.parent.unlink(l)
	return nil
}

func (l *logger) obj() *Obj {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	res := &Obj{logger: l}
	for _, mw := range l.config.Pre {
		res = mw(l.config, res)
	}
	return res
}

func (l *logger) Dict() *Obj {
	return l.obj().Dict()
}
func (l *logger) Array() *Obj {
	return l.obj().Array()
}
func (l *logger) Str(v string) *Obj {
	return l.obj().Str(v)
}
func (l *logger) Fmt(t string, vs ...any) *Obj {
	return l.obj().Str(fmt.Sprintf(t, vs...))
}
func (l *logger) Float(v float64) *Obj {
	return l.obj().Float(v)
}
func (l *logger) Int(v int) *Obj {
	return l.obj().Int(v)
}
func (l *logger) Bool(v bool) *Obj {
	return l.obj().Bool(v)
}
func (l *logger) Null() *Obj {
	return l.obj().Null()
}
func (l *logger) Bytes(d []byte) *Obj {
	return l.obj().Bytes(d)
}

func Match(pat string) (map[string]map[string]int, error) {
	return root.Match(pat)
}

func (l *logger) Match(pat string) (map[string]map[string]int, error) {
	re, err := regexp.Compile(pat)
	if err != nil {
		return nil, err
	}
	return l.match(nil, re), nil
}

func (l *logger) match(dst map[string]map[string]int, pat *regexp.Regexp) map[string]map[string]int {
	l.mu.Lock()
	defer l.mu.Unlock()
	if dst == nil {
		dst = map[string]map[string]int{}
	}
	for k, v := range l.config.Labels {
		if !pat.MatchString(k) {
			continue
		}
		pMap := dst[l.config.pkg]
		if pMap == nil {
			pMap = map[string]int{}
			dst[l.config.pkg] = pMap
		}
		pMap[l.config.Localize(k)] = v
	}
	for k := range l.children {
		k.match(dst, pat)
	}
	return dst
}
